package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/serbanmunteanu/xm-golang-task/config"
)

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(config config.KafkaConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Addr},
		Topic:   config.Topic,
		GroupID: config.ReaderGroupID,
	})

	return &Consumer{Reader: reader}
}
