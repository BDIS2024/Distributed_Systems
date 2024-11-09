package main

import (
	"Assignment_4/proto"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewMutualExclusionClient(conn)

	stream, err := client.RequestAccess(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = stream.Send(&proto.Request{NodeId: "1", Timestamp: 1})
	if err != nil {
		log.Fatal(err)
	}

	// Handle the stream here (e.g., receive responses)
	// Example: response, err := stream.Recv()

	// Close the stream if necessary
	// Example: stream.CloseSend()
}
