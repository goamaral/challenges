package server

import (
	"esl-challenge/api/gen/userpb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

const ServiceName = "esl-challenge"

type Server struct {
	*grpc.Server
	healthServer *health.Server
}

func NewServer() (*Server, error) {
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			logrus.Panic(p)
			return status.Error(codes.Unknown, "panic")
		}),
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New())),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logrus.New())),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)

	userpb.RegisterUserServiceServer(grpcServer, NewUserServiceServer())

	healthServer := health.NewServer()
	healthServer.SetServingStatus(UserServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	return &Server{
		Server:       grpcServer,
		healthServer: healthServer,
	}, nil
}

func (s Server) SetServingStatus(service string, servingStatus grpc_health_v1.HealthCheckResponse_ServingStatus) {
	s.healthServer.SetServingStatus(service, servingStatus)
}
