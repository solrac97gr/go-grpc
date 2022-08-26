package main

import (
	"github.com/solrac97gr/go-grpc/database"
	"github.com/solrac97gr/go-grpc/server"
	"github.com/solrac97gr/go-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	list, err := net.Listen("tcp", ":5070")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository(
		"postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable",
	)

	serverTest := server.NewTestServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	testpb.RegisterTestServiceServer(s, serverTest)
	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
