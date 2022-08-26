package main

import (
	"context"
	"github.com/solrac97gr/go-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}

	defer conn.Close()

	client := testpb.NewTestServiceClient(conn)

	DoUnary(client)
	DoClientStreaming(client)
	DoServerStreaming(client)
	DoBidirectionalStreaming(client)
}

func DoUnary(client testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "t1",
	}

	res, err := client.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("No se pudo obtener la respuesta: %v", err)
	}

	log.Printf("response from server: %v", res)
}

func DoClientStreaming(client testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q8t1",
			Answer:   "45",
			Question: "9x5",
			TestId:   "t1",
		},
		{
			Id:       "q9t1",
			Answer:   "100",
			Question: "10x10",
			TestId:   "t1",
		},
		{
			Id:       "q10t1",
			Answer:   "9",
			Question: "3x3",
			TestId:   "t1",
		},
	}
	stream, err := client.SetQuestions(context.Background())
	if err != nil {
		log.Fatalf("Error while calling SetQuestions: %v", err)
	}

	for _, question := range questions {
		stream.Send(question)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciving response: %v", err)
	}
	log.Printf("Response from server: %v", msg)
}

func DoServerStreaming(client testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{TestId: "t1"}

	stream, err := client.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GetStudentsPerTest : %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error sending the Student: %v", err)
			break
		}
		log.Printf("Response form the server: %v", msg)
	}

}

func DoBidirectionalStreaming(client testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer: "45",
	}

	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := client.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("Error while calling TakeTest: %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			stream.Send(&answer)
		}

	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while reading stream: %v", err)
				break
			}
			log.Printf("Response form the server: %v", res)
		}
		close(waitChannel)
	}()
	<-waitChannel
}
