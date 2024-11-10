package main

import (
	proto "Assignment_4-real-real/proto"
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
var name = ""

func main() {
	//logs
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	//connection
	reader := bufio.NewReader(os.Stdin)

	//stream

	node1stream, node1conn := connectToHost("5050")
	node2stream, node2conn := connectToHost("5051")

	defer node1conn.Close()
	defer node2conn.Close()

	broadcast("i want to join", []proto.DmutexService_DmutexClient{node1stream, node2stream})

	fmt.Println("Enter your name:")
	name, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	name = strings.TrimSpace(name)


	log.Printf("Client sent request: Name: %s, Message: %s, Timestamp: (%d)\n", msg.Name, msg.Message, counter)
	waitc := make(chan bool)
	//donec := make(chan bool)

	//go retrieveMessage(waitc, donec, stream)
	//go sendMessage(donec, stream, msg.Name)

	<-waitc

}

func connectToHost(host string) (proto.DmutexService_DmutexClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient("localhost:"+host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}

	client := proto.NewDmutexServiceClient(conn)

	stream, err := client.Dmutex(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	return stream, conn
}

func broadcast(message string, nodes []proto.DmutexServiceClient) {
	msg := proto.Req{
		Name:      name,
		Message:   message,
		Timestamp: counter,
	}

	for _, node := range nodes {
		err := node.Send(&msg)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func retrieveMessage(waitc chan bool, donec chan bool, stream proto.DmutexService_DmutexClient) {
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
				log.Fatal("Error receiving message:", err)
				waitc <- true
				return
			}
			counter = max(counter, in.Timestamp) + 1

			log.Printf("Client recieved response: Name: %s, Message: %s, Timestamp: (%d) at: %d\n", in.Name, in.Message, in.Timestamp, counter)
		}
	}
}

func sendMessage(donec chan bool, stream proto.DmutexService_DmutexClient, username string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		message = strings.TrimSpace(message)

		if message == "leave" {
			err = stream.Send(&proto.Req{
				Name:      username,
				Message:   "has left the chat.",
				Timestamp: counter,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
			err = stream.CloseSend()
			if err != nil {
				log.Fatal("Error closing stream:", err)
			}

			donec <- true
			return
		}
		counter++

		msg := proto.Req{
			Name:      username,
			Message:   message,
			Timestamp: counter,
		}
		err = stream.Send(&msg)
		if err != nil {
			log.Fatal("Error sending message:", err)
		}
		log.Printf("Client sent request: Name: %s, Message: %s, Timestamp: (%d)\n", msg.Name, msg.Message, counter)
	}
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}
