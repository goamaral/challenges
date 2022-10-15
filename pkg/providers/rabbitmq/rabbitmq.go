package rabbitmq

import (
	"context"
	"esl-challenge/pkg/env"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitmqProvider interface {
	Close()
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
}

type rabbitmqProvider struct {
	*amqp.Channel
	conn *amqp.Connection
}

func NewRabbitmqProvider() (_ RabbitmqProvider, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to rabbitmq
	var conn *amqp.Connection
	for {
		conn, err = amqp.Dial(env.GetOrDefault("RABBIT_MQ_URL", "amqp://localhost:5672"))
		if err != nil {
			select {
			case <-time.After(time.Second):
				logrus.Info("Waiting for rabbitmq to be ready")
			case <-ctx.Done():
				return nil, err
			}
		} else {
			break
		}
	}

	// Open channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &rabbitmqProvider{conn: conn, Channel: channel}, nil
}

func (p *rabbitmqProvider) Close() {
	p.Channel.Close()
	p.conn.Close()
}
