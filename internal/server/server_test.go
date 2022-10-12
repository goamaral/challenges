package server_test

import (
	"esl-challenge/internal/repository"
	"esl-challenge/internal/server"
	"net"
	"testing"
)

func initServer(t *testing.T, userRepository repository.UserRepository) (net.Listener, *server.Server) {
	grpcServer, err := server.NewServer(userRepository)
	if err != nil {
		t.Fatal(err)
	}
	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatal(err)
	}
	return lis, grpcServer
}
