package rabbitmq_ext

import (
	"challenge/pkg/env"
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Channel struct {
	*amqp.Channel
	conn *amqp.Connection
}

func NewChannel() (Channel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to rabbitmq
	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial(env.GetOrDefault("RABBIT_MQ_URL", "amqp://localhost:5672"))
		if err != nil {
			select {
			case <-time.After(time.Second):
				logrus.Info("Waiting for rabbitmq to be ready")
			case <-ctx.Done():
				return Channel{}, err
			}
		} else {
			break
		}
	}

	// Open channel
	channel, err := conn.Channel()
	if err != nil {
		return Channel{}, conn.Close()
	}

	return Channel{Channel: channel, conn: conn}, nil
}

func (p *Channel) Close() {
	p.Channel.Close()
	p.conn.Close()
}
