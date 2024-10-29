package main

import (
	pb "Lecture_5_Exercise/proto"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error")
	}

	client := pb.NewTimeServiceClient(conn)

	time, err := client.GetTime(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("error")
	}
	fmt.Println(time)
}
