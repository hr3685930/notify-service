package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronjan/hunch"
	"github.com/golang-module/carbon"
	"github.com/streadway/amqp"
	"notify-service/pkg/queue"
	"reflect"
	"strings"
	"time"
)

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	QueueName    string
	Exchange     string
	Key          string
	DLKQueueName string
	DLKExchange  string
	DLKKey       string
	MQUrl        string
	Prefix       string
}

func NewRabbitMQ(user, pass, host, port, vhost, prefix string) queue.Queue {
	mqUrl := "amqp://" + user + ":" + pass + "@" + host + ":" + port + vhost
	return &RabbitMQ{MQUrl: mqUrl, Prefix: prefix}
}

func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(r.MQUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("ampq connect error %s", err))
	}
	channel, err := conn.Channel()
	if err != nil {
		return errors.New(fmt.Sprintf("ampq channel error %s", err))
	}

	r.conn = conn
	r.channel = channel
	return nil
}

func (r *RabbitMQ) ProducerConnect() queue.Queue {
	channel, err := r.conn.Channel()
	if err != nil {
		amqperr := err.(*amqp.Error)
		if amqperr.Code == amqp.ChannelError {
			_ = r.Connect()
			return &RabbitMQ{MQUrl: r.MQUrl, channel: r.channel, conn: r.conn, Prefix: r.Prefix}
		}
		panic(fmt.Sprintf("ampq channel error %s", err))
	}

	r.channel = channel
	return &RabbitMQ{MQUrl: r.MQUrl, channel: channel, conn: r.conn, Prefix: r.Prefix}
}

func (r *RabbitMQ) ConsumerConnect() queue.Queue {
	channel, err := r.conn.Channel()
	if err != nil {
		amqperr := err.(*amqp.Error)
		if amqperr.Code == amqp.ChannelError {
			_ = r.Connect()
			return &RabbitMQ{MQUrl: r.MQUrl, channel: r.channel, conn: r.conn, Prefix: r.Prefix}
		}
		panic(fmt.Sprintf("ampq channel error %s", err))
	}

	r.channel = channel
	return &RabbitMQ{MQUrl: r.MQUrl, channel: channel, conn: r.conn, Prefix: r.Prefix}
}

func (r *RabbitMQ) Close() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func (r *RabbitMQ) Topic(topic string) {
	r.SetExchange(topic)
}

func (r *RabbitMQ) SetExchange(exchange string) {
	r.Exchange = exchange
}

func (r *RabbitMQ) SetQueue(queue string) {
	comma := strings.Index(queue, ".")
	r.QueueName = r.Prefix + "_" + strings.ToLower(strings.Replace(queue[comma+1:], ".", "_", -1))
}

func (r *RabbitMQ) SetKey(key string) {
	comma := strings.Index(key, ".")
	r.Key = strings.ToLower(key[comma+1:])
}

func (r *RabbitMQ) SetDLKExchange() {
	r.DLKExchange = "delay." + r.Exchange
}

func (r *RabbitMQ) SetDLKQueue(sleep int32) {
	r.DLKQueueName = fmt.Sprintf("delay-%d_%s", sleep, r.QueueName)
}

func (r *RabbitMQ) SetDLKKey(sleep int32) {
	r.DLKKey = fmt.Sprintf("delay-%d.%s", sleep, r.Key)
}

func (r *RabbitMQ) Producer(job queue.JobBase, delay int32) {
	if r.Exchange == "" {
		panic("not exchange")
	}

	queueName := reflect.TypeOf(job).String()
	r.SetKey(queueName[1:] + "." + r.Prefix)

	message, err := json.Marshal(job)
	r.failOnErr(err, "Umarshal failed")

	//create Exchange
	err = r.channel.ExchangeDeclare(r.Exchange, "topic", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	//pub
	err = r.channel.Publish(r.Exchange, r.Key, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(message),
			DeliveryMode: 2,
			Timestamp:    time.Now(),
			Headers: map[string]interface{}{
				"delay": delay,
			},
		})
	_ = r.channel.Close()
}

//DLK
func (r *RabbitMQ) DLK(base queue.JobBase, sleep int32, headers map[string]interface{}) {
	//create DLK Exchange
	if r.DLKExchange == "" {
		r.SetDLKExchange()
		err := r.channel.ExchangeDeclare(r.DLKExchange, "topic", true, false, false, false, nil)
		r.failOnErr(err, "Failed to declare an delay_exchange")
	}

	//create DLK queue
	if r.DLKQueueName == "" {
		r.SetDLKQueue(sleep)
		r.SetDLKKey(sleep)
		args := make(amqp.Table)
		args["x-dead-letter-exchange"] = r.Exchange
		args["x-dead-letter-routing-key"] = r.Key
		args["x-message-ttl"] = 1000 * sleep
		q, err := r.channel.QueueDeclare(r.DLKQueueName, true, false, false, false, args)
		r.failOnErr(err, "Failed to declare a delay_queue")
		//bind DLK
		err = r.channel.QueueBind(q.Name, r.DLKKey, r.DLKExchange, false, nil)
		r.failOnErr(err, "Failed to declare a delay_queue bind")
	}

	message, err := json.Marshal(base)
	r.failOnErr(err, "Umarshal failed")
	err = r.channel.Publish(r.DLKExchange, r.DLKKey, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(message),
			DeliveryMode: 1,
			Timestamp:    time.Now(),
			Headers:      headers,
		})

}

// Consumer  单进程来保证顺序消费
func (r *RabbitMQ) Consumer(base queue.JobBase, sleep, retry, timeout int32) {
	queueName := reflect.TypeOf(base).String()
	r.SetQueue(queueName[1:])
	r.SetKey(queueName[1:] + ".*")
	// create Exchange
	err := r.channel.ExchangeDeclare(r.Exchange, "topic", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an delay_exchange")

	// create Queue
	q, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare a queue")

	//bind
	err = r.channel.QueueBind(q.Name, r.Key, r.Exchange, false, nil)
	r.failOnErr(err, "Failed to declare a queue bind")

	//Qos
	err = r.channel.Qos(1, 0, false)
	//Consumer
	messages, err := r.channel.Consume(q.Name, "", false, false, false, false, nil)
	forever := make(chan bool)

	go func() {
		for d := range messages {
			err := json.Unmarshal(d.Body, base)
			if err != nil {
				r.failOnErr(err, "Unmarshal error")
			}

			// delay
			if pDelay, ok := d.Headers["delay"].(int32); ok && pDelay > 0 {
				r.DLK(base, pDelay, nil)
				err = d.Ack(false)
				if err != nil {
					r.ExportErr(queue.Err(err), d)
				}
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			if timeout == 0 {
				ctx = context.Background()
			}

			// retry and delay
			_, err = hunch.Retry(ctx, int(retry)+1, func(ctx context.Context) (interface{}, error) {
				err = base.Handler()
				if err != nil {
					time.Sleep(time.Second * time.Duration(sleep))
				}
				return nil, err
			})
			if err != nil {
				r.ExportErr(queue.Err(err), d)
			}

			err = d.Ack(false)
			if err != nil {
				r.ExportErr(queue.Err(err), d)
			}
			cancel()
		}
	}()

	<-forever
}

func (r *RabbitMQ) ExportErr(err error, d amqp.Delivery) {
	e := err.(*queue.Error)
	r.Err(queue.FailedJobs{
		Connection: "rabbitmq",
		Queue:      r.QueueName,
		Message:    string(d.Body),
		Exception:  err.Error(),
		Stack:      e.Stack(),
		FiledAt:    carbon.Now(),
	})
}

func (r *RabbitMQ) Err(failed queue.FailedJobs) {
	queue.ErrJob <- failed
}
