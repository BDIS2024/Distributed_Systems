package main

import (
	pb "Lecture_5_Exercise/proto"
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTimeServiceServer
}

func (s *server) GetTime(ctx context.Context, in *pb.Empty) (*pb.Time, error) {
	tid := time.Now()
	t := tid.String()
	return &pb.Time{Time: t}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("error")
	}

	pb.RegisterTimeServiceServer(grpcServer, &server{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("error")
	}
}
