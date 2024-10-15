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

	// setup cli with send message command, get messages command
	// send messages takes 1 arg the messages to send
	// get messages takes 0 args

	// *to send message
	// call sendmessage with message arg
	// client code will compute a lamport timestamp (nodenr, eventnr)
	// pass message and timstamp into sendmessage remote function
	// to compute nodenr when client joining send join request, server will compute an id for client, when new client joins compute another id and send back

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
