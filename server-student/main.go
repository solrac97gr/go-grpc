package main

import (
	"github.com/solrac97gr/go-grpc/database"
	"github.com/solrac97gr/go-grpc/server"
	"github.com/solrac97gr/go-grpc/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	list, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository(
		"postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable",
	)

	serverStudent := server.NewStudentServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, serverStudent)
	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
