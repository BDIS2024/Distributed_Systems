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
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name string
var auctionservers []string
var auctionserverconnections []proto.AuctionServiceClient
var output []*proto.Outcome

// fix locks
var lock sync.Mutex

func main() {
	start := time.Now()
	auctionservers = append(auctionservers, "5050", "5051", "5052")
	connectToServers()

	name = getName()

	wait := make(chan bool)

	go prompt(wait)

	for {
		lock.Lock()
		output = nil
		for i := 0; i < len(auctionserverconnections); i++ {
			result, err := auctionserverconnections[i].Result(context.Background(), &proto.Empty{})
			if err != nil {
				fmt.Printf("Auctionserver %s is down.\n", auctionservers[i])
				continue
			}
			output = append(output, result)
		}
		lock.Unlock()
		if len(output) == 0 {
			fmt.Println("All auction servers are down.")
			wait <- true
			break
		}
		if !ongoing(output) {
			fmt.Println("Acution has ended.")
			fmt.Printf("The highest bidder was %s with a bid of %d.\n", output[0].HighestBidder, output[0].HighestBid)
			wait <- true
			break
		}
	}

	<-wait
	elapsed := time.Since(start)
	fmt.Println("Program done.")
	fmt.Printf("Time taken: %s\n", elapsed.String())
}

func ongoing(output []*proto.Outcome) bool {
	lock.Lock()
	for _, outcome := range output {
		if outcome.Status == "Ongoing" {
			lock.Unlock()
			return true
		}
	}
	lock.Unlock()
	return false
}

func prompt(stop chan bool) {
	inputchannel := make(chan string)

	go input(inputchannel, stop)

	for {
		select {
		case <-stop:
			stop <- true
			return
		case input := <-inputchannel:
			switch {
			case strings.HasPrefix(input, "bid"):
				var acks []*proto.Ack
				for _, auctionserver := range auctionserverconnections {
					bid(auctionserver, input, &acks)
				}
				if len(acks) > 0 {
					printSpaces()
					fmt.Printf("Bid with %s was: %s\n", input, acks[0].Outcome)
				}
			case input == "result":
				for _, auctionserver := range auctionserverconnections {
					result(auctionserver)
				}
			}
		}
	}
}

func input(inputchannel chan string, stop chan bool) {
	for {
		select {
		case <-stop:
			stop <- true
			return
		default:
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("To bid type 'Bid <amount>' and press enter.")
			fmt.Println("To get the status of the auction type 'Result' and press enter.")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			input = strings.TrimSpace(input)
			input = strings.ToLower(input)
			inputchannel <- input
			time.Sleep(1 * time.Second)
		}
	}
}

func bid(client proto.AuctionServiceClient, bid string, acks *[]*proto.Ack) {
	amountstr := strings.Split(bid, " ")[1]
	amountint, err := strconv.ParseInt(amountstr, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := client.Bid(context.Background(), &proto.Amount{Bid: amountint, Bidder: name})
	if err != nil {
		log.Fatalln(err)

	}
	*acks = append(*acks, result)
}

func result(client proto.AuctionServiceClient) {
	result, err := client.Result(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalln(err)
	}
	printSpaces()
	if result.Status == "Ongoing" {
		fmt.Printf("The highest bid is %d by %s.\n", result.HighestBid, result.HighestBidder)
	} else {
		fmt.Println("Auction has ended.")
		fmt.Printf("The highest bid was %d by %s.\n", result.HighestBid, result.HighestBidder)
	}
}

func printSpaces() {
	for i := 0; i < 100; i++ {
		fmt.Println()
	}
}

func connectToServers() {
	for _, server := range auctionservers {
		hostname := "localhost:" + server
		conn, err := grpc.NewClient(hostname, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err.Error())
		}
		client := proto.NewAuctionServiceClient(conn)
		auctionserverconnections = append(auctionserverconnections, client)
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

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
