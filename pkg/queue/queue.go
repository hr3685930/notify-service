package queue

import "runtime/debug"

type JobBase interface {
	Handler() (queueErr *Error)
}

type Queue interface {
	// Connect connect
	Connect() error
	// ProducerConnect Producer
	ProducerConnect() Queue
	// ConsumerConnect consumer connect
	ConsumerConnect() Queue
	// Topic abstract topic
	Topic(topic string)
	// Producer delay
	Producer(job JobBase, delay int32)
	// Consumer sleep retry
	Consumer(job JobBase, sleep, retry, timeout int32)
	// Err report
	Err(failed FailedJobs)

	Close()
}

type Error struct {
	s     string
	stack string
}

func (qe *Error) Error() string {
	return qe.s
}

func Err(err error) *Error {
	return &Error{
		s:     err.Error(),
		stack: string(debug.Stack()),
	}
}

func (qe *Error) Stack() string {
	return qe.stack
}
