package main

import (
	proto "Assignment_4/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type dmutexserver struct {
	proto.UnimplementedDMutexServiceServer
}

func (s *dmutexserver) DistributedMutexService(ctx context.Context, in *proto.Req) (*proto.Resp, error) {
	return &proto.Resp{Name: "node", Timestamp: 1}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	proto.RegisterDMutexServiceServer(grpcServer, &dmutexserver{})

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server started on :5050")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
