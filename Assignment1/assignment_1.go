package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	channel := make(chan bool)
	var wg sync.WaitGroup

	for {
		for i := 1; i < 6; i++ {
			wg.Add(1)
			go philosopher(i, channel, &wg)
		}

		for i := 1; i < 6; i++ {
			message := <-channel
			fmt.Println(message)
		}

		wg.Wait()

		time.Sleep(1 * time.Second)
		fmt.Println()
	}
}

func philosopher(id int, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	philo := phil{eating: false, thinking: false, id: id}
	ch <- true
	fmt.Println(philo.id, "done", philo.eating, philo.thinking)

}

type phil struct {
	eating   bool
	thinking bool
	id       int
}
