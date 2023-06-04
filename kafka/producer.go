package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/serbanmunteanu/xm-golang-task/config"
)

type Producer interface {
	Produce(key []byte, value []byte) error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(config config.KafkaConfig) Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Addr),
		Topic:                  config.Topic,
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
		Async:                  true,
	}

	return &producer{writer: writer}
}

func (p *producer) Produce(key []byte, value []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})
}
