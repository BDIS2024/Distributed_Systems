package main

import (
	pb "Lecture_5/proto"
	"context"
	"time"
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

}
