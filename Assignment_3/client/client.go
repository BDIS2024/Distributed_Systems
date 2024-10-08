package main

import (
	proto "Assignment_3/proto"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connection error")
	}

	client := proto.ChittyChatServiceClient(conn)

	time, err := client.GetMessages(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("error2")
	}
	fmt.Println(time)
}
