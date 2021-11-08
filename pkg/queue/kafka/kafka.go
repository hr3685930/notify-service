package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/aaronjan/hunch"
	"github.com/golang-module/carbon"
	"github.com/rfyiamcool/go-timewheel"
	"notify-service/pkg/queue"
	"reflect"
	"strings"
	"time"
)

type Kafka struct {
	cli            sarama.Client
	Brokers        []string
	ConsumerTopics []string
	ProducerTopic  string
	Prefix         string
}

func NewKafka(urls, prefix string) queue.Queue {
	brokers := strings.Split(urls, ",")
	return &Kafka{Prefix: prefix, Brokers: brokers}
}

type consumerGroupHandler struct {
	k         *Kafka
	Sleep     int32
	Retry     int32
	TimeOut   int32
	Message   queue.JobBase
	TimeWheel *timewheel.TimeWheel
}

func (k *Kafka) Connect() error {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	client, err := sarama.NewClient(k.Brokers, config)
	if err != nil {
		return err
	}
	k.cli = client
	return nil
}

func (k *Kafka) ProducerConnect() queue.Queue {
	return &Kafka{cli: k.cli, Prefix: k.Prefix, ProducerTopic: k.ProducerTopic, ConsumerTopics: k.ConsumerTopics, Brokers: k.Brokers}
}

func (k *Kafka) ConsumerConnect() queue.Queue {
	return &Kafka{cli: k.cli, Prefix: k.Prefix, ProducerTopic: k.ProducerTopic, ConsumerTopics: k.ConsumerTopics, Brokers: k.Brokers}
}

func (k *Kafka) Topic(topic string) {
	k.ProducerTopic = topic
	k.ConsumerTopics = []string{topic}
}

func (k *Kafka) failOnErr(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func (k *Kafka) Producer(job queue.JobBase, delay int32) {
	p, err := sarama.NewSyncProducerFromClient(k.cli)
	k.failOnErr(err, "producer create producer error")

	jName := reflect.TypeOf(job).String()
	comma := strings.Index(jName, ".")
	// 增加key,hash到同一partition保证顺序消费,但发生rebalance时也不能保证顺序性
	// 避免发生rebalance 1.不允许临时增加组下消费者 2.不允许更改partition数
	key := strings.ToLower(strings.Replace(jName[comma+1:], ".", "_", -1))

	var headers []sarama.RecordHeader
	header := sarama.RecordHeader{
		Key:   []byte("delay"),
		Value: queue.Int32ToBytes(delay),
	}

	headers = append(headers, header)

	msg := &sarama.ProducerMessage{
		Topic:   k.ProducerTopic,
		Key:     sarama.StringEncoder(key),
		Headers: headers,
	}

	message, err := json.Marshal(job)
	k.failOnErr(err, "Umarshal failed")
	msg.Value = sarama.ByteEncoder(message)
	_, _, err = p.SendMessage(msg)
	k.failOnErr(err, "send error")
	_ = p.Close()
}

func (k *Kafka) Consumer(job queue.JobBase, sleep, retry, timeout int32) {
	jName := reflect.TypeOf(job).String()
	comma := strings.Index(jName, ".")
	groupID := strings.ToLower(strings.Replace(jName[comma+1:]+"."+k.Prefix, ".", "_", -1))
	group, err := sarama.NewConsumerGroupFromClient(groupID, k.cli)
	k.failOnErr(err, "Consumer group err")
	ctx := context.Background()
	for { //防止rebalance后结束
		topics := k.ConsumerTopics
		handler := &consumerGroupHandler{k: k, Message: job, Retry: retry, Sleep: sleep, TimeOut: timeout}
		err = group.Consume(ctx, topics, handler)
		k.failOnErr(err, "Consumer err")
	}
}

func (k *Kafka) Err(failed queue.FailedJobs) {
	queue.ErrJob <- failed
}

func (k *Kafka) Close() {
	_ = k.cli.Close()
}

func (k *Kafka) ExportErr(err error, msg, groupID string) {
	e := err.(*queue.Error)
	k.Err(queue.FailedJobs{
		Connection: "kafka",
		Queue:      groupID + "_" + k.Prefix,
		Message:    msg,
		Exception:  err.Error(),
		Stack:      e.Stack(),
		FiledAt:    carbon.Now(),
	})
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	tw, _ := timewheel.NewTimeWheel(1*time.Second, 360, timewheel.SetSyncPool(true))
	c.TimeWheel = tw
	c.TimeWheel.Start()
	return nil
}
func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	c.TimeWheel.Stop()
	return nil
}
func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := json.Unmarshal(msg.Value, c.Message)
		if err != nil {
			sess.MarkMessage(msg, "")
		}
		jName := reflect.TypeOf(c.Message).String()
		comma := strings.Index(jName, ".")
		if strings.ToLower(jName[comma+1:]) == string(msg.Key) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeOut)*time.Second)
			if c.TimeOut == 0 {
				ctx = context.Background()
			}

			headers := make(map[string]interface{}, 1)
			for _, value := range msg.Headers {
				headers[string(value.Key)] = queue.BytesToInt32(value.Value)
			}
			delay := headers["delay"].(int32)
			if delay > 0 {
				jsonRes := msg.Value
				// interface copy
				msgHandler := reflect.New(reflect.ValueOf(c.Message).Elem().Type()).Interface().(queue.JobBase)
				_ = c.TimeWheel.Add(time.Duration(delay)*time.Second, func() {
					_, err = hunch.Retry(ctx, int(c.Retry)+1, func(ctx context.Context) (interface{}, error) {
						err := json.Unmarshal(jsonRes, &msgHandler)
						if err != nil {
							sess.MarkMessage(msg, "")
						}
						err = msgHandler.Handler()
						if err != nil {
							c.TimeWheel.Sleep(time.Duration(c.Sleep) * time.Second)
						}
						return nil, err
					})
					if err != nil {
						c.k.ExportErr(queue.Err(err), string(jsonRes), string(msg.Key))
					}
					sess.MarkMessage(msg, "")
				})
				cancel()
				continue
			}

			_, err = hunch.Retry(ctx, int(c.Retry)+1, func(ctx context.Context) (interface{}, error) {
				err = c.Message.Handler()
				if err != nil {
					c.TimeWheel.Sleep(time.Duration(c.Sleep) * time.Second)
				}
				return nil, err
			})

			if err != nil {
				fmt.Print(err.Error())
				c.k.ExportErr(queue.Err(err), string(msg.Value), string(msg.Key))
			}
			cancel()
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
