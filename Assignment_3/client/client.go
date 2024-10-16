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

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("type join to join a chat session, type exit to exit program")

		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		message = strings.TrimSpace(message)

		if message == "join" {
			fmt.Println("joined")
			waitc := make(chan struct{})
			go retrive(waitc, stream)

			go send(waitc, stream)

			<-waitc
		}
		if message == "exit" {
			break
		}

	}

}

func retrive(waitc chan struct{}, stream proto.ChittyChatService_ChatServiceClient) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			close(waitc)
			return
		}

		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("message: %s timestamp: %d\n", in.Message, in.Timestamp)
	}
}

func send(waitc chan struct{}, stream proto.ChittyChatService_ChatServiceClient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("enter message")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}

		message = strings.TrimSpace(message)

		if message == "leave" {
			fmt.Println("!!!!!")
			close(waitc)
			return
		}
		stream.Send(&proto.Message{Message: message, Timestamp: 1})
	}
}
