package queue

import (
	"notify-service/pkg/db"
	"sync"
)

var QueueStore sync.Map

var MQ Queue

var ErrJob chan FailedJobs

func GetQueueDrive(c string) Queue {
	v, ok := QueueStore.Load(c)
	if ok {
		return v.(Queue)
	}
	return nil
}

func NewConsumer(topic string, job JobBase, sleep, retry, timeout int32) {
	AutoMigrate()
	mq := MQ.ConsumerConnect()
	mq.Topic(topic)
	mq.Consumer(job, sleep, retry, timeout)
}

func NewProducer(topic string, job JobBase, delay int32) {
	mq := MQ.ProducerConnect()
	mq.Topic(topic)
	mq.Producer(job, delay)
}

func AutoMigrate() {
	_ = db.Orm.AutoMigrate(&FailedJobs{})
	ErrJob = make(chan FailedJobs, 1)
	go func() {
		for {
			select {
			case failedJob := <-ErrJob:
				db.Orm.Save(&failedJob)
			}
		}
	}()
}
