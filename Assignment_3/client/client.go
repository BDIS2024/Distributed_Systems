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
		log.Fatal(err.Error())
	}

	client := proto.NewChittyChatServiceClient(conn)

	message, errr := client.SendMessage(context.Background(), &proto.Message{Message: "hej", Timestamp: 1})
	if errr != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Sent message: %s, with timestamp: %d\n", message.Message, message.Timestamp)

	messages, err := client.GetMessages(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, message := range messages.Messages {
		fmt.Printf(" - %s, : %d \n", message.Message, message.Timestamp)
	}
}
