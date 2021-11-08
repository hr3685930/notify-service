package configs

import "notify-service/pkg/config"

// Queue default once
type Queue struct {
	RabbitMQ RabbitMQ
	Kafka    Kafka
}

type RabbitMQ struct {
	config.RabbitMQDrive
	Host     string `default:"127.0.0.1"`
	Port     string `default:"5672"`
	VHost    string `default:"/"`
	Username string `default:"admin"`
	Password string `default:"admin"`
}

type Kafka struct {
	config.KafkaDrive
	Addr   string `default:"127.0.0.1:9092"`
}



