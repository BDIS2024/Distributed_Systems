package main

import (
	proto "Assignment_4/proto"
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

var counter int32
var name string
var port string
var state = "RELEASED"

type Node struct {
	stream proto.DmutexService_DmutexClient
	conn   *grpc.ClientConn
	port   string
}

func main() {
	//logs
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// set the node name
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter node identifier:")
	name, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	name = strings.TrimSpace(name)

	// set node server port
	fmt.Println("Enter server port:")
	port, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	port = strings.TrimSpace(port)

	// get other node streams
	node1 := connectToHost("5051")
	node2 := connectToHost("5052")

	defer node1.conn.Close()
	defer node2.conn.Close()

	message := "i want to join"

	// broadcast to other nodes
	log.Printf("Client: %s, sent message: %s, with timestamp: %d to :%s and :%s\n", name, message, counter, node1.port, node2.port)
	broadcast(message, []proto.DmutexService_DmutexClient{node1.stream, node2.stream})

	waitc := make(chan bool)
	donec := make(chan bool)

	go retrieveMessage(waitc, donec)
	//go sendMessage(donec, stream, msg.Name)

	<-waitc

}

func connectToHost(host string) Node {
	conn, err := grpc.NewClient("localhost:"+host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}

	client := proto.NewDmutexServiceClient(conn)

	stream, err := client.Dmutex(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	node := Node{
		stream: stream,
		conn:   conn,
		port:   host,
	}

	return node
}

func broadcast(message string, nodes []proto.DmutexService_DmutexClient) {
	msg := proto.Message{
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

func retrieveMessage(waitc chan bool, donec chan bool) {
	self := connectToHost(port)

	for {
		select {
		case <-donec:
			waitc <- true
			return
		default:
			in, err := self.stream.Recv()
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
			err = stream.Send(&proto.Message{
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

		msg := proto.Message{
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
