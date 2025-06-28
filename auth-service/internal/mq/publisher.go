package mq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitPublisher struct {
	channel *amqp.Channel
}

func NewRabbitPublisher(rabbitURL string) (*RabbitPublisher, error) {
	const maxRetries = 5
	const retryDelay = 5 * time.Second

	var conn *amqp.Connection
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logger.Log.Info("connecting to RabbitMQ", zap.String("url", rabbitURL), zap.Int("attempt", attempt))

		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			// ok connected
			break
		}

		logger.Log.Warn("failed to connect to RabbitMQ, will retry", zap.Error(err))
		time.Sleep(retryDelay)
	}

	if err != nil {
		logger.Log.Error("giving up connecting to RabbitMQ", zap.Error(err))
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Log.Error("failed to open a channel", zap.Error(err))
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"user.events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		logger.Log.Error("failed to declare exchange", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("RabbitMQ publisher connected successfully")
	return &RabbitPublisher{channel: ch}, nil
}

func (p *RabbitPublisher) PublishUserSignedUp(event UserSignUpEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"user.events", // exchange
		"user.signup", // routing key
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *RabbitPublisher) Close() {
	if err := p.channel.Close(); err != nil {
		log.Printf("failed to close RabbitMQ channel: %v", err)
	} else {
		log.Println("RabbitMQ channel closed successfully")
	}
}
