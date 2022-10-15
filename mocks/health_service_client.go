// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"

	mock "github.com/stretchr/testify/mock"
)

// HealthServiceClient is an autogenerated mock type for the HealthServiceClient type
type HealthServiceClient struct {
	mock.Mock
}

// Check provides a mock function with given fields: ctx, in, opts
func (_m *HealthServiceClient) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption) (*grpc_health_v1.HealthCheckResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *grpc_health_v1.HealthCheckResponse
	if rf, ok := ret.Get(0).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) *grpc_health_v1.HealthCheckResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc_health_v1.HealthCheckResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Watch provides a mock function with given fields: ctx, in, opts
func (_m *HealthServiceClient) Watch(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption) (grpc_health_v1.Health_WatchClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 grpc_health_v1.Health_WatchClient
	if rf, ok := ret.Get(0).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) grpc_health_v1.Health_WatchClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(grpc_health_v1.Health_WatchClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpc_health_v1.HealthCheckRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewHealthServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewHealthServiceClient creates a new instance of HealthServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHealthServiceClient(t mockConstructorTestingTNewHealthServiceClient) *HealthServiceClient {
	mock := &HealthServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}