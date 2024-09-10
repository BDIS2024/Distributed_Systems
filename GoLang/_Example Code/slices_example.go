package main

import "fmt"

func main() {

	fmt.Println("REFERENCE EXAMPLE")
	var x string
	x = "ciao"
	fmt.Println(x)
	fmt.Println(&x)

	fmt.Println()
	fmt.Println("SLICE EXAMPLES")
	// init 1
	var numbers []int
	numbers = []int{1, 2, 3, 5, 4}
	fmt.Println(numbers)
	fmt.Println(len(numbers))
	fmt.Println(cap(numbers))

	// init 2
	var numbers2 = []int{99, 100, 90, 80, 70}

	// init 3
	numbers3 := []int{9, 7, 8}

	// init 4, with allocation
	var numbers4 = make([]int, 2, 10) // int, length 2, capacity 2
	fmt.Println(numbers4)

	numbers4[1] = 5 // assign 5 to the 2nd element
	// numbers4[2] = 3 // ERROR: assign 3 to the 3rd element
	fmt.Println(numbers4)
	fmt.Println(len(numbers4))
	fmt.Println(cap(numbers4))

	numbers4 = []int{20, 21, 23}
	fmt.Println(len(numbers4))
	fmt.Println(cap(numbers4))

	fmt.Println()
	fmt.Println()
	fmt.Println(cap(numbers))
	numbers = append(numbers, 58)
	fmt.Println(cap(numbers))

	numbers = append(numbers, 42)
	fmt.Println(cap(numbers))

	fmt.Println()
	fmt.Println()

	fmt.Println(numbers)
	fmt.Println(numbers2)
	fmt.Println(numbers3)
	fmt.Println(numbers4)
}
