package main

import (
	proto "Assignment_4/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	stream proto.DmutexService_DmutexClient
	conn   *grpc.ClientConn
	port   string
}

var clientNodePair Node

var knwonNodesNode []Node
var knownNodes []string

var replies int

func main() {
	//logs
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// setup peer-to-peer network
	knownNodes = append(knownNodes, "5050")
	knownNodes = append(knownNodes, "5051")
	knownNodes = append(knownNodes, "5052")

	// get pair port and connect
	port := getPort()
	clientNodePair = connectToHost(port)
	connectToPair()

	sendToAllNodes()

	// get other node streams
	//node1 := connectToHost("5051")
	//node2 := connectToHost("5052")
	//
	//defer node1.conn.Close()
	//defer node2.conn.Close()
	//
	//message := "i want to join"
	//
	//// broadcast to other nodes
	//log.Printf("Client: %s, sent message: %s, with timestamp: %d to :%s and :%s\n", name, message, counter, node1.port, node2.port)
	//broadcast(message, []proto.DmutexService_DmutexClient{node1.stream, node2.stream})

	waitc := make(chan bool)
	//donec := make(chan bool)

	//go retrieveMessage(waitc, donec)
	//go sendMessage(donec, stream, msg.Name)

	<-waitc

}

func getPort() string {
	var port string
	var err error

	// Port
	if len(os.Args) > 1 {

		fmt.Printf("test:%v\n", os.Args[1])
		port = os.Args[1]
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		fmt.Println("Enter port number:")
		reader := bufio.NewReader(os.Stdin)
		port, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	port = strings.TrimSpace(port)
	return port
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

func connectToPair() {
	msg := proto.Message{
		Name:      "",
		Message:   "Connect",
		Timestamp: 0,
	}
	clientNodePair.stream.Send(&msg)

	for i := 0; i < len(knownNodes); i++ {
		if knownNodes[i] == clientNodePair.port {
			knownNodes = remove(knownNodes, i)
			break
		}
	}
}

func getTime() int32 {
	// unimplemented
	return 0
}

func sendToAllNodes() {
	// init nodes
	if len(knownNodes) > len(knwonNodesNode) {
		for i := 0; i < len(knownNodes); i++ {
			knwonNodesNode = append(knwonNodesNode, connectToHost(knownNodes[i]))
		}
	}

	msg := proto.Message{
		Name:      "",
		Message:   "Request",
		Timestamp: getTime(),
	}

	// send
	for i := 0; i < len(knwonNodesNode); i++ {
		err := knwonNodesNode[i].stream.Send(&msg)
		if err != nil {
			log.Fatal("Error sending message:", err)
		}
	}
}

/*
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

		log.Printf("Client sent request: Name: %s, Message: %s, Timestamp: (%d)\n", msg.Name, msg.Message, counter)
	}
}*/
/*
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

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
} */

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
