package main

import (
	"fmt"
	"time"
)

type TCPPackage struct {
	message string
	seq     int
	ack     int
}

func main() {
	channel := make(chan TCPPackage)

	go client(channel)
	go server(channel)
	fmt.Println("Client and server spawned")
	time.Sleep(3 * time.Second)
}

func client(channel chan TCPPackage) {
	message := TCPPackage{"request", 0, 0}
	channel <- message
	fmt.Printf("Client sent request package with sequence: %d\n", message.seq)

	response := <-channel
	if response.seq == 0 {
		channel <- response
		fmt.Printf("Client pulled out again too early with sequence : %d\n", response.seq)
		fmt.Printf("Client sent back with sequence: %d\n", response.seq)
	} else {
		fmt.Printf("Client retrieved response with sequence: %d\n", response.seq)
		fmt.Printf("Client retrieved response with ack: %d\n", response.ack)
		response.seq += 1
		response.ack += 1
		channel <- response
	}

}

func server(channel chan TCPPackage) {
	message := <-channel
	fmt.Printf("Server retrieved package with sequence: %d\n", message.seq)
	message.seq += 1
	message.ack = 1

	channel <- message
	fmt.Printf("Server sent back response with seq: %d and ack: %d\n", message.seq, message.ack)

	response := <-channel
	if response.seq == 1 {
		channel <- response
	} else {
		fmt.Printf("Server retrieved with sequence: %d, ack: %d and message: %s\n", response.seq, response.ack, response.message)
	}

}

/*
1)[Easy] Implement the TCP/IP Handshake using threads. This is not realistic (since the protocol should run across a network)
but your implementation needs to show that you have a good understanding of the protocol.
*/
