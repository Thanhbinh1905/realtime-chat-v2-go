package mq

import (
	"encoding/json"
	"log"

	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitPublisher struct {
	channel *amqp.Channel
}

func NewRabbitPublisher(rabbitURL string) (*RabbitPublisher, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		logger.Log.Error("failed to connect to RabbitMQ", zap.Error(err))
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
		return nil, err
	}

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
