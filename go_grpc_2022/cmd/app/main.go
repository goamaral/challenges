package main

import (
	"challenge/internal/repository"
	"challenge/internal/server"
	"challenge/internal/service"
	"challenge/pkg/env"
	"challenge/pkg/providers/postgres"
	"challenge/pkg/providers/rabbitmq"
	"fmt"
	"net"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load env file
	if env.GetOrDefault("ENV", "development") == "development" {
		err := godotenv.Load()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to load .env")
		}
	}

	// Providers
	postgresProvider, err := postgres.NewPostgresProvider()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize postgres provider")
	}
	rabbitmqProvider, err := rabbitmq.NewRabbitmqProvider()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize rabbitmq provider")
	}
	defer rabbitmqProvider.Close()

	// Declare rabbitmq exchange
	err = rabbitmqProvider.ExchangeDeclare(service.RabbitmqChangesTopic, "topic", true, false, false, false, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to declare rabbitmq changes exchange")
	}

	// Repositories
	userRepo := repository.NewUserRepository(postgresProvider)

	// Services
	rabbitmqSvc := service.NewRabbitmqService(rabbitmqProvider)

	// Server
	grpcServer, err := server.NewServer(userRepo, rabbitmqSvc)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize grpc server")
		return
	}

	port := 3000
	logrus.Infof("Starting grpc server on port %d", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on port %d", port)
		return
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start grpc server")
		return
	}
}
