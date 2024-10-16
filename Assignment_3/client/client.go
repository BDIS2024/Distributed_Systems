package main

import (
	proto "Assignment_3/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}

	client := proto.NewChittyChatServiceClient(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("enter command\n join to join a session\nquit to quit program\n")
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read")
		}

		command = strings.TrimSpace(command)
		if command == "join" {
			fmt.Println("joined")
			for {
				fmt.Print("enter command\n leave to leave\n send to send message\n get to get messages\n")
				todo, erro := reader.ReadString('\n')
				todo = strings.TrimSpace(todo)
				if erro != nil {
					log.Fatalf("failed to read")
				}

				if todo == "leave" {
					fmt.Println("left")
					break
				}

				if todo == "send" {
					fmt.Println("enter message")
					tosend, erro := reader.ReadString('\n')
					if erro != nil {
						log.Fatalf("failed to read")
					}
					tosend = strings.TrimSpace(tosend)
					sendMessage(client, tosend)
					getMessages(client)
				}

				if todo == "get" {
					getMessages(client)
				}
			}
		}
		if command == "quit" {
			break
		}

	}
	fmt.Println("done")
	// setup cli with send message command, get messages command, join command, quit command,
	// send messages takes 1 arg the messages to send
	// get messages takes 0 args
	// join takes 0 args
	// quit take 0 args

	// *to send message
	// call sendmessage with message arg
	// client code will compute a lamport timestamp (nodenr, eventnr)
	// pass message and timstamp into sendmessage remote function
	// to compute nodenr when client joining send join request, server will compute an id for client, when new client joins compute another id and send back

	// this is not live chatting
	// can use bi directional streaming*

}

func sendMessage(client proto.ChittyChatServiceClient, messag string) {
	message, err := client.SendMessage(context.Background(), &proto.Message{Message: messag, Timestamp: 1})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Sent message: %s, with timestamp: %d\n", message.Message, message.Timestamp)
}

func getMessages(client proto.ChittyChatServiceClient) {
	messages, err := client.GetMessages(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(messages.Messages[len(messages.Messages)-1].Message)
}
