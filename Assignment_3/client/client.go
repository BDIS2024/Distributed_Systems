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

var counter int32 = 0

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	client := proto.NewChittyChatServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your username")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	username = strings.TrimSpace(username)

	for {
		fmt.Println("Type 'join' to join a chat session.")
		fmt.Println("Type 'exit' to exit the program")

		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		message = strings.TrimSpace(message)

		if message == "join" {
			fmt.Printf("%s has joined the chat. (%d)\n", username, counter)

			stream, err := client.ChatService(context.Background())
			if err != nil {
				log.Fatal(err.Error())
			}

			err = stream.Send(&proto.ClientMessage{
				Name:      username,
				Message:   "has joined the chat.",
				Timestamp: counter,
			})
			if err != nil {
				log.Println(err.Error())
			}

			waitc := make(chan bool)
			donec := make(chan bool)

			go retrieveMessage(waitc, donec, stream)
			go sendMessage(donec, stream, username)

			<-waitc
		} else if message == "exit" {
			fmt.Println("Exiting program...")
			break
		}
	}
}

func retrieveMessage(waitc chan bool, donec chan bool, stream proto.ChittyChatService_ChatServiceClient) {
	for {
		select {
		case <-donec:
			waitc <- true
			return
		default:
			in, err := stream.Recv()
			if err == io.EOF {
				waitc <- true
				return
			}
			if err != nil {
				log.Println("Error receiving message:", err)
				waitc <- true
				return
			}
			counter = max(counter, in.Timestamp) + 1
			fmt.Printf("%s : %s (%d)\n", in.Name, in.Message, counter)
		}
	}
}

func sendMessage(donec chan bool, stream proto.ChittyChatService_ChatServiceClient, username string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		message = strings.TrimSpace(message)

		if message == "leave" {
			fmt.Printf("%s has left the chat. (%d)\n", username, counter)

			err = stream.Send(&proto.ClientMessage{
				Name:      username,
				Message:   "has left the chat.",
				Timestamp: counter,
			})
			if err != nil {
				log.Println(err.Error())
			}

			err = stream.CloseSend()
			if err != nil {
				log.Println("Error closing stream:", err)
			}

			donec <- true
			return
		}
		counter++
		err = stream.Send(&proto.ClientMessage{
			Name:      username,
			Message:   message,
			Timestamp: counter,
		})
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}
