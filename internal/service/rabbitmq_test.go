package service_test

import (
	"context"
	"encoding/json"
	"esl-challenge/internal/service"
	"esl-challenge/mocks"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRabbitmqService_PublishChanges(t *testing.T) {
	entity := 0

	rabbitmqProvider := mocks.NewRabbitmqProvider(t)
	rabbitmqProvider.On("PublishWithContext", mock.Anything, service.RabbitmqChangesTopic, "user.create", false, false, mock.Anything).
		Return(func(_ context.Context, _ string, _ string, _ bool, _ bool, msg amqp.Publishing) error {
			assert.Equal(t, msg.ContentType, "application/json")
			var res int
			err := json.Unmarshal(msg.Body, &res)
			if assert.NoError(t, err) {
				assert.Equal(t, entity, res)
			}

			return nil
		})

	s := service.NewRabbitmqService(rabbitmqProvider)
	err := s.PublishChanges(context.Background(), entity, service.EntityType_USER, service.Action_CREATE)
	assert.NoError(t, err)
}
