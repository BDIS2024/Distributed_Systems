package main

import (
	proto "Assignment_5/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	start := time.Now()
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}

	client := proto.NewAuctionServiceClient(conn)

	wait := make(chan bool)

	go prompt()

	for {
		result, err := client.Result(context.Background(), &proto.Empty{})
		if err != nil {
			log.Fatalln(err)
		}
		if result.Status != "Ongoing" {
			fmt.Println("Acution has ended.")
			fmt.Printf("The highest bidder was %s with a bid of %d.\n", result.HighestBidder, result.HighestBid)
			wait <- true
			break
		}
	}

	<-wait
	elapsed := time.Since(start)
	fmt.Println("Program done.")
	fmt.Printf("Time taken: %s\n", elapsed.String())
}

func prompt() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your bid: ")
		bid, err := reader.ReadString('\n')
		bid = strings.TrimSpace(bid)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(bid)
	}
}
