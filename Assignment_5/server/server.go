package main

import (
	proto "Assignment_5/proto"
	"context"
	"fmt"
	"log"
	"net"

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
	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServiceServer(grpcServer, &AuctionServer{})

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalln(err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Server started listening of port :5050")
}
