package main

import (
	"esl-challenge/internal/repository"
	"esl-challenge/internal/server"
	"esl-challenge/pkg/env"
	"esl-challenge/pkg/providers/postgres"
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

	// Repositories
	userRepository := repository.NewUserRepository(postgresProvider)

	// Server
	grpcServer, err := server.NewServer(userRepository)
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
