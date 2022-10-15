// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	amqp091 "github.com/rabbitmq/amqp091-go"

	mock "github.com/stretchr/testify/mock"
)

// RabbitmqProvider is an autogenerated mock type for the RabbitmqProvider type
type RabbitmqProvider struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *RabbitmqProvider) Close() {
	_m.Called()
}

// ExchangeDeclare provides a mock function with given fields: name, kind, durable, autoDelete, internal, noWait, args
func (_m *RabbitmqProvider) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp091.Table) error {
	ret := _m.Called(name, kind, durable, autoDelete, internal, noWait, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool, bool, bool, bool, amqp091.Table) error); ok {
		r0 = rf(name, kind, durable, autoDelete, internal, noWait, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PublishWithContext provides a mock function with given fields: ctx, exchange, key, mandatory, immediate, msg
func (_m *RabbitmqProvider) PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp091.Publishing) error {
	ret := _m.Called(ctx, exchange, key, mandatory, immediate, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool, bool, amqp091.Publishing) error); ok {
		r0 = rf(ctx, exchange, key, mandatory, immediate, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRabbitmqProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewRabbitmqProvider creates a new instance of RabbitmqProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRabbitmqProvider(t mockConstructorTestingTNewRabbitmqProvider) *RabbitmqProvider {
	mock := &RabbitmqProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
