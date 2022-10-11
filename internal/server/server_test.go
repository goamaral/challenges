package server_test

import (
	"esl-challenge/internal/server"
	"net"
	"testing"
)

func initServer(t *testing.T) (net.Listener, *server.Server) {
	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatal(err)
	}
	grpcServer, err := server.NewServer()
	if err != nil {
		t.Fatal(err)
	}
	return lis, grpcServer
}
