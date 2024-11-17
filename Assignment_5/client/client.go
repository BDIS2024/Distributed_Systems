package main

import (
	proto "Assignment_5/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name string

func main() {
	start := time.Now()
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}

	client := proto.NewAuctionServiceClient(conn)
	name = getName()

	wait := make(chan bool)

	go prompt(client)

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

func prompt(client proto.AuctionServiceClient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("To bid type 'Bid <amount>' and press enter.")
		fmt.Println("To get the status of the auction type 'Result' and press enter.")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		input = strings.TrimSpace(input)

		switch {
		case strings.HasPrefix(input, "Bid"):
			bid(client, input)
		case input == "Result":
			result(client)
		}

		fmt.Println(bid)
	}
}

func bid(client proto.AuctionServiceClient, bid string) {
	amountstr := strings.Split(bid, " ")[1]
	amountint, err := strconv.ParseInt(amountstr, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := client.Bid(context.Background(), &proto.Amount{Bid: amountint, Bidder: name})
}

func result(client proto.AuctionServiceClient) {
	result, err := client.Result(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalln(err)
	}
}

func getName() string {
	var name string
	var err error

	if len(os.Args) > 1 {

		name = os.Args[1]

	} else {
		fmt.Println("Enter your name:")
		reader := bufio.NewReader(os.Stdin)
		name, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	name = strings.TrimSpace(name)
	return name
}
