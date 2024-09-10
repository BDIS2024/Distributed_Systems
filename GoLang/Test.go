package main

// Sollution INPIRED and NOT COPIED by:
// https://medium.com/@bararviv0120/dining-philosophers-problem-golang-6b6a1cc2fdb7

/*

	This sollution is based on randomness to avoid a deadlock.
	Every philosopher will try to grab and then reserve their left fork, then their right.
	If at any point they fail to do so they will give up their forks and think for a random amount of time.
	This ensures that at some point every philosopher will get a window to grab and reserve their forks.

	Observed runtime.
	We have recorded some runtimes to shape the readers expectations.

	1. 46s
	2. 17s
	3. 43s
	4. 29s

*/

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	forks := [5]fork{}

	for i := 0; i < len(forks); i++ {
		forks[i] = createFork(i)
	}

	phils := [5]phil{}

	for i := 0; i < len(phils); i++ {
		phils[i] = createPhil(i, &forks[i], &forks[(i+1)%5])
	}

	for i := 0; i < len(phils); i++ {
		go forkRoutine(forks[i])
	}

	for i := 0; i < len(phils); i++ {
		go philRoutine(&phils[i])
	}

	start := time.Now()

	for {
		all := true

		for i := 0; i < len(phils); i++ {
			if phils[i].eatCount < 3 {
				all = false
				break
			}
		}

		if all {
			fmt.Println("All Philosophers are now full :)\nTime taken: ", time.Since(start))
			break
		}
	}
}

type phil struct {
	id       int
	eatCount int
	leftPtF  chan int
	leftFtP  chan int
	rightPtF chan int
	rightFtP chan int
}

type fork struct {
	id       int
	used     bool
	leftPtF  chan int
	leftFtP  chan int
	rightPtF chan int
	rightFtP chan int
}

func createFork(idA int) fork {
	return fork{
		id:       idA,
		used:     false,
		leftFtP:  make(chan int, 1),
		rightFtP: make(chan int, 1),
		leftPtF:  nil,
		rightPtF: nil,
	}
}

func createPhil(idA int, leftFork *fork, rightFork *fork) phil {
	p := phil{
		id:       idA,
		eatCount: 0,
		leftFtP:  leftFork.leftFtP,
		rightFtP: rightFork.rightFtP,
		leftPtF:  make(chan int, 1),
		rightPtF: make(chan int, 1),
	}
	leftFork.leftPtF = p.leftPtF
	rightFork.rightPtF = p.rightPtF
	return p
}

/*
	Phil to Fork Communication (PtF)

	1 - Request Fork
	2 - Clear Usage of Fork

*/

/*
	Fork to Phil Communication (FtP)

	1 - Fork Request Denied
	2 - Fork Request Accepted

*/

func forkRoutine(f fork) {
	for {
		checkForkChannel(&f, false)
		checkForkChannel(&f, true)
		//time.Sleep(1 * time.Second)
	}
}

func checkForkChannel(f *fork, dir bool) {
	// Step 1: Get direction
	var in chan int
	var out chan int
	if dir {
		in = f.leftPtF
		out = f.leftFtP
	} else {
		in = f.rightPtF
		out = f.rightFtP
	}

	// Step 2: Check if there are requests
	//fmt.Printf("Id: %d, len: %d in:", f.id, len(in))
	//fmt.Print(in, "\n")
	if len(in) == 0 {
		return
	}

	// Step 3: Read and answer request
	result := <-in

	switch result {
	case 1: // Fork is requested
		if f.used {
			out <- 1 // Denied
		} else {
			out <- 2 // Accepted
			f.used = true
		}
		break
	case 2: // Clear fork
		f.used = false
		break
	default:
		fmt.Printf("ERROR: %d NOT RECOGNIZED FOR checkForkChannel!\n", result)
		break
	}
}

func philRoutine(p *phil) {
	for {

		// Step 1: Request left fork
		if len(p.leftPtF) == 1 {
			think(p, 2)
			continue
		}

		p.leftPtF <- 1 // Send
		//fmt.Println("sent to: ", p.leftPtF)
		resultL := <-p.leftFtP

		if resultL == 1 { // fork is denied try again
			think(p, 2)
			continue
		}
		fmt.Printf("I %d got my left fork!\n", p.id)

		// Step 2: Request right fork
		if len(p.rightPtF) == 1 {
			think(p, 2)
			continue
		}
		p.rightPtF <- 1 // Send
		resultR := <-p.rightFtP

		if resultR == 1 { // fork is denied try again
			p.rightPtF <- 2 // free the left fork
			think(p, 2)
			continue
		}
		fmt.Printf("I %d got my right fork!\n", p.id)

		// Step 3: Sucess

		eat(p)
		fmt.Printf("I %d ate!(%d/3)\n", p.id, p.eatCount+1)

		// free fork resources
		p.leftPtF <- 2
		p.rightPtF <- 2

		// increment eat and announce
		p.eatCount += 1
		think(p, 7)
	}

}

func eat(p *phil) { // If we keep recursively sending
	think(p, 1)
}

func think(p *phil, mult int64) { // If we keep recursively sending
	fmt.Printf("I %d am thinking!\n", p.id)
	time.Sleep(time.Duration(rand.Int63n(1e9) * mult))
}
