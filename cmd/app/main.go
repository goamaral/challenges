package main

import (
	"esl-challenge/internal/server"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

func main() {
	port := 3000

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on port %d", port)
		return
	}

	grpcServer, err := server.NewServer()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize grpc server")
		return
	}

	logrus.Infof("Starting grpc server on port %d", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start grpc server")
		return
	}
}
