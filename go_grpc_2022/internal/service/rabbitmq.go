package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"

	"challenge/api/gen/userpb"
	"challenge/pkg/rabbitmq_ext"
)

type Topic string

const (
	Topic_USERS Topic = "users"
)

type Subject string

const (
	Subject_CREATED Subject = "created"
	Subject_UPDATED Subject = "updated"
	Subject_DELETED Subject = "deleted"
)

type RabbitmqService interface {
	Publish(ctx context.Context, event proto.Message) error
}

type rabbitmqService struct {
	channel rabbitmq_ext.Channel
}

func NewRabbitmqService(channel rabbitmq_ext.Channel) (RabbitmqService, error) {
	err := channel.ExchangeDeclare(string(Topic_USERS), "topic", true, false, false, false, nil)
	if err != nil {
		return rabbitmqService{}, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return rabbitmqService{channel: channel}, nil
}

func (s rabbitmqService) Publish(ctx context.Context, event proto.Message) error {
	topic, subject, ok := getTopicAndSubjectFromEvent(event)
	if ok {
		return fmt.Errorf("failed to get topic and subject from %T", event)
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := amqp.Publishing{ContentType: "application/json", Body: body}
	err = s.channel.PublishWithContext(ctx, string(topic), string(subject), false, false, msg)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func getTopicAndSubjectFromEvent(event proto.Message) (Topic, Subject, bool) {
	switch event.(type) {
	case *userpb.Event_UserCreated:
		return Topic_USERS, Subject_CREATED, true
	default:
		return "", "", false
	}
}
