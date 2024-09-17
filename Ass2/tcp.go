package main

type TCPPackage struct{
	message string
	ack int
	seq int
}

func main() {
	channel := make(chan TCPPackage)

	go client(channel)
	go server(channel)
}

func client(chan cqwe){

}

func server(chan wqe){

}
