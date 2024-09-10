package main

import (
	"fmt"
	"os"
)

func main() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] { // iterates through slice from 1 to its end
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}
