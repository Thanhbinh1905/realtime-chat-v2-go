package event

import (
	"context"

	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"
	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokerAddress string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) Publish(ctx context.Context, key string, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		logger.Log.Error("Failed to publish message to Kafka", zap.Error(err), zap.String("key", key))
	}
	return err
}
func (p *KafkaProducer) Close() error {
	if err := p.writer.Close(); err != nil {
		logger.Log.Error("Failed to close Kafka writer", zap.Error(err))
		return err
	}
	return nil
}
