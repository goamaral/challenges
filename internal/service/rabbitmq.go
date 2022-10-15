package service

import (
	"context"
	"encoding/json"
	"esl-challenge/pkg/providers/rabbitmq"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Action string

const (
	Action_CREATE Action = "create"
	Action_UPDATE Action = "update"
	Action_DELETE Action = "delete"
)

type EntityType string

const (
	EntityType_USER EntityType = "user"
)

const RabbitmqChangesTopic = "topic_changes"

type RabbitmqService interface {
	PublishChanges(ctx context.Context, entity any, entityName EntityType, action Action) error
}

type rabbitmqService struct {
	provider rabbitmq.RabbitmqProvider
}

func NewRabbitmqService(provider rabbitmq.RabbitmqProvider) RabbitmqService {
	return &rabbitmqService{provider: provider}
}

func (s rabbitmqService) PublishChanges(ctx context.Context, entity any, entityType EntityType, action Action) error {
	// Marshall entity
	body, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	// Publish message in changes topic exchange
	msg := amqp.Publishing{ContentType: "application/json", Body: body}
	err = s.provider.PublishWithContext(ctx, RabbitmqChangesTopic, fmt.Sprintf("%s.%s", entityType, action), false, false, msg)
	if err != nil {
		return err
	}

	return nil
}
