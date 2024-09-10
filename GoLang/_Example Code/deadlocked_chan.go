package main

import "fmt"

var ch1 = make(chan int)
var ch2 = make(chan string)

func thread1() {
	fmt.Println("I'm thread1 and I received the following int:", <-ch2)
	ch1 <- 5
}

func thread2() {
	fmt.Println("I'm thread2 and I received the following string:", <-ch1)
	ch2 <- "Hello"

}

func main() {
	go thread1()
	go thread2()

	for {
	}

}
