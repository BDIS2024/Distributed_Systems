package main

import (
	proto "Assignment_5/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServiceServer
}

func (s *AuctionServer) Bid(context.Context, *proto.Amount) (*proto.Ack, error) {
	log.Println("Bid")
	return &proto.Ack{Outcome: "Success"}, nil
}

func (s *AuctionServer) Result(context.Context, *proto.Empty) (*proto.Outcome, error) {
	log.Println("Result")
	return &proto.Outcome{HighestBid: "500", HighestBidder: "Nicky"}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServiceServer(grpcServer, &AuctionServer{})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
}
