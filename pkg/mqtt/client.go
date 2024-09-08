package mqtt

import (
	"bank-authentication-system/pkg/config"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type MQTT struct {
	Channel *amqp091.Channel
	Queue   string
}

// NewConnection makes a connection to rabbit cluster and returns the MQTT type.
func NewConnection(cfg config.MQTT) (*MQTT, error) {
	var mq MQTT

	conn, err := amqp091.Dial(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("rabbitMQ connection failed: %w", err)
	}

	mq.Channel, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("getting rabbitMQ channel failed: %w", err)
	}

	_, err = mq.Channel.QueueDeclare(
		cfg.Queue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("creating queue failed: %w", err)
	}

	mq.Queue = cfg.Queue

	return &mq, nil
}
