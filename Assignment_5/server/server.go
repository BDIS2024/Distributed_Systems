package main

import (
	proto "Assignment_5/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServiceServer
}

func (s *AuctionServer) Bid(context.Context, *proto.Amount) (*proto.Ack, error) {
	fmt.Println("Bid")
	return &proto.Ack{Outcome: "Success"}, nil
}

func (s *AuctionServer) Result(context.Context, *proto.Empty) (*proto.Outcome, error) {
	fmt.Println("Result")
	return &proto.Outcome{HighestBid: 500, HighestBidder: "Nicky", Status: "Ongoing"}, nil
}

func main() {
	// Setup Server

	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServiceServer(grpcServer, &AuctionServer{})

	port := getPort()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Server started listening on port %s\n", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	var port string
	var err error

	// Port
	if len(os.Args) > 1 {

		port = os.Args[1]

	} else {
		fmt.Println("Enter port number:")
		reader := bufio.NewReader(os.Stdin)
		port, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	port = strings.TrimSpace(port)
	port = ":" + port
	return port
}

// TODO: LEADER ELECTION & REPLICATION
