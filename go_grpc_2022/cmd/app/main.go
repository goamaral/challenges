package main

import (
	"fmt"
	"net"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"challenge/internal/repository"
	"challenge/internal/server"
	"challenge/internal/service"
	"challenge/pkg/gorm_ext"
	"challenge/pkg/rabbitmq_ext"
)

func main() {
	// Load env file
	godotenv.Load()

	// Providers
	db, err := gorm_ext.ConnectToDatabase(
		gorm_ext.NewPostgresDialector(gorm_ext.NewPostgresDSN()),
		false,
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to db")
	}
	rabbitChan, err := rabbitmq_ext.NewChannel()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize rabbitmq channel")
	}
	defer rabbitChan.Close()

	// Services
	rabbitmqSvc, err := service.NewRabbitmqService(rabbitChan)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize to rabbitmq service")
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)

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
