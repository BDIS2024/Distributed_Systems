package main

import (
	proto "Assignment_4-real-real/proto"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

func main() {
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	grpcServer := grpc.NewServer()
	proto.RegisterDmutexServiceServer(grpcServer, &DmutexServer{})

	fmt.Println("Enter port number:")
	reader := bufio.NewReader(os.Stdin)
	port, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	port = strings.TrimSpace(port)
	port = ":" + port

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started on %s", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

type DmutexServer struct {
	proto.UnimplementedDmutexServiceServer
}

type MessageObject struct {
	ClientName string
	Message    string
	Timestamp  int32
}

type MessageHandler struct {
	Clients map[string]proto.DmutexService_DmutexServer
	Lock    sync.Mutex
}

var handler = MessageHandler{
	Clients: make(map[string]proto.DmutexService_DmutexServer),
}

var counter int32 = 0

func (s *DmutexServer) Dmutex(stream proto.DmutexService_DmutexServer) error {
	errorChan := make(chan error)

	go retrieveMessagesFromClient(stream, errorChan)

	return <-errorChan
}

func retrieveMessagesFromClient(stream proto.DmutexService_DmutexServer, errorChan chan error) {

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			errorChan <- err
			return
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			errorChan <- err
			return
		}

		// if len(message.Message) > 128 {
		// 	sendErrorToCLient(clientName, "Message has to be under 128 characters.")
		// 	continue
		// }

		counter = max(counter, message.Timestamp) + 1
		fmt.Println(message)
		log.Printf("Server recieved request: Name: %s, Message: %s, Timestamp: (%d) at %d\n", message.Name, message.Message, message.Timestamp, counter)
		//broadcastMessageToClients(message)
	}
}

func broadcastMessageToClients(message *proto.Ack) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()
	counter++
	for clientName, clientStream := range handler.Clients {
		err := clientStream.Send(&proto.Ack{
			Name:      message.Name,
			Message:   message.Message,
			Timestamp: counter,
		})
		if err != nil {
			log.Printf("Error sending message to %s: %v", clientName, err)
		}
	}
	log.Printf("Server sent response: Name: %s, Message: %s, Timestamp: (%d)\n", message.Name, message.Message, counter)
}

func sendErrorToCLient(clientName string, erro string) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()
	counter++

	err := handler.Clients[clientName].Send(&proto.Ack{
		Name:      "Server",
		Message:   erro,
		Timestamp: counter,
	})
	if err != nil {
		log.Printf("Error sending message to %s: %v", clientName, err)

	}
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}
