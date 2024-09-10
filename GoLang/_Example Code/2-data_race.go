package main

import (
	"fmt"
	"sync"
	"time"
)

const N = 10000000

var balance = 0

var arbiter sync.Mutex

func worker() {
	for i := 0; i < N; i++ {
		arbiter.Lock()
		balance++
		// arbiter.Unlock()
	}
	fmt.Println("Done")
}

func main() {

	go worker()
	go worker()

	for {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(balance)
	}
}
