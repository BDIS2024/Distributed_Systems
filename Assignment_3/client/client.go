package main

import (
	proto "Assignment_3/proto"
	"bufio"
	"context"
	"fmt"
	"io"
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

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your username")

	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	username = strings.TrimSpace(username)

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("type join to join a chat session, type exit to exit program")

		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		message = strings.TrimSpace(message)

		if message == "join" {
			fmt.Println("Joining chat...")
			err = stream.Send(&proto.ClientMessage{
				Name:      username,
				Message:   "has joined the chat.",
				Timestamp: "1",
			})
			if err != nil {
				log.Println(err.Error())
			}

			waitc := make(chan struct{})

			go retrieveMessage(waitc, stream)
			go sendMessage(waitc, stream, username)

			<-waitc
		}
		if message == "exit" {
			break
		}

	}

}

func retrieveMessage(waitc chan struct{}, stream proto.ChittyChatService_ChatServiceClient) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			close(waitc)
			return
		}

		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%s : %s (%s)\n", in.Name, in.Message, in.Timestamp)
	}
}

func sendMessage(waitc chan struct{}, stream proto.ChittyChatService_ChatServiceClient, username string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}

		message = strings.TrimSpace(message)

		if message == "leave" {
			fmt.Println("Leaving chat...")
			err = stream.Send(&proto.ClientMessage{
				Name:      username,
				Message:   "has left the chat.",
				Timestamp: "1",
			})
			if err != nil {
				log.Println(err.Error())
			}

			close(waitc)

			err = stream.CloseSend()
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
		err = stream.Send(&proto.ClientMessage{
			Name:      username,
			Message:   message,
			Timestamp: "1",
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
