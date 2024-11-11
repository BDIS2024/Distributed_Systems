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

var counter int32 = 0
var port string

type Node struct {
	stream proto.DmutexService_DmutexClient
	conn   *grpc.ClientConn
	port   string
}

var clientNodePair Node

var knwonNodesNode []Node
var knownNodes []string

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
	port = getPort()
	clientNodePair = connectToHost(port)
	connectToPair()

	replies = 0
	hasEnoughReplies = make(chan bool, 1)
	requestingCriticalSection = make(chan bool, 1)

	go replyRoutine()

	askForCriticalSection()

	waitc := make(chan bool)

	<-waitc

}

var hasEnoughReplies chan bool
var requestingCriticalSection chan bool

var replies int
var requestTimeStamp int32

func isRequestingCriticalSection() bool {
	return len(requestingCriticalSection) >= 1
}

func askForCriticalSection() {

	sendRequestToAllNodes()

	requestingCriticalSection <- true

	<-hasEnoughReplies

	// access critical section
	fmt.Printf("%s ACCESSED THE CRITICAL SECTION AT LAMPORT TIMESTAMP: %d\n", port, counter)
	log.Printf("%s ACCESSED THE CRITICAL SECTION %d\n", port, counter)

	// free the other routine
	replyToStoredReplies()
	<-requestingCriticalSection
	hasEnoughReplies <- true

}

func replyRoutine() {
	for {
		message, err := clientNodePair.stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		if message.Name == port {
			log.Printf("Error received message from self: %v", message)
			continue
		}

		fmt.Printf("I recived message: %v\n", message)
		log.Printf("%s recived message: %v\n", port, message)

		var recievedTimestamp = message.Timestamp
		counter = max(counter, recievedTimestamp)
		counter = counter + 1

		if message.Message == "Reply" {
			if isRequestingCriticalSection() {
				//storedReplies = append(storedReplies, message.Name) Remove // if not working
				replies++
				if replies >= len(knownNodes) { // check if we got enough replies
					hasEnoughReplies <- true
					replies = 0
					<-hasEnoughReplies // "lock" this method temporarily while we acces critical section so we dont send stored replies and then at the same time add another element to the array
				}

			} else { // we should not be getting this message
				log.Printf("Error recived unexpected message: %v", message)
				continue
			}
		} else if message.Message == "Request" {
			// Main logic for algorythm
			if isRequestingCriticalSection() {

				// determine who gets prio

				if message.Timestamp == requestTimeStamp { // If timestamps are equal determaine by port number
					if message.Name > port { //
						// The other port has prio
						sendReply(message.Name)
					} else {
						// We have prio
						storedReplies = append(storedReplies, message.Name)
					}
				} else if message.Timestamp > requestTimeStamp { //
					// The other port has prio
					sendReply(message.Name)
				} else {
					// We have prio
					storedReplies = append(storedReplies, message.Name)
				}

			} else {
				sendReply(message.Name)
			}

		} else { // we should not be getting this message
			log.Printf("Error recived unknown message: %v\n", message)
			continue
		}

	}
}

var storedReplies []string

func replyToStoredReplies() {

	fmt.Printf("Replying to stored replies (%v)\n", len(storedReplies))
	log.Printf("%s Replying to stored replies (%v) at %d\n", port, len(storedReplies), counter)

	for i := 0; i < len(storedReplies); i++ {
		sendReply(storedReplies[i])
	}

	storedReplies = []string{}
}

func sendReply(reciverPort string) {
	counter = counter + 1
	for i := 0; i < len(knwonNodesNode); i++ {
		if knwonNodesNode[i].port == reciverPort { // send message

			msg := proto.Message{
				Name:      port,
				Message:   "Reply",
				Timestamp: counter,
			}

			err := knwonNodesNode[i].stream.Send(&msg)
			if err != nil {
				log.Fatal("Error sending message:", err)
			} else {
				fmt.Printf("I sent message: Name:%v message:%v Timestamp: %v\n", msg.Name, msg.Message, msg.Timestamp)
				return
			}
		}
	}

	// could not find port
	log.Printf("Error: Could not send message to port %v, port was not found\n", port)
	fmt.Printf("Error: Could not send message to port %v, port was not found\n", port)
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

func sendRequestToAllNodes() {

	// init nodes
	if len(knownNodes) > len(knwonNodesNode) {
		for i := 0; i < len(knownNodes); i++ {
			knwonNodesNode = append(knwonNodesNode, connectToHost(knownNodes[i]))
		}
	}

	counter = counter + 1
	requestTimeStamp = counter

	log.Printf("%v sending request at %v", port, requestTimeStamp)

	msg := proto.Message{
		Name:      port,
		Message:   "Request",
		Timestamp: requestTimeStamp,
	}

	// send
	for i := 0; i < len(knwonNodesNode); i++ {
		err := knwonNodesNode[i].stream.Send(&msg)
		if err != nil {
			log.Fatal("Error sending message:", err)
		}
	}
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}


// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
