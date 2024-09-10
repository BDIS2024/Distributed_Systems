package main

import (
	"fmt"
)

func main() {
	var x, y int
	fmt.Scanln(&x, &y)

	if x >= y {
		t := x
		x = y
		y = t
	}

	fmt.Printf("%d %d", x, y)
}
