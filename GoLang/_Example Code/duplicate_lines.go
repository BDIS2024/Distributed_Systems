package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int) // this is a standard map, from string (keys) to int (values)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++ // it checks entry 'input.Text()' and increments its value
		// 	new keys initialised to 0
	}
	// NOTE: ignoring potential errors from input.Err()
	for line, n := range counts {
		if n > 1 { // no paranthesis around condition, as with for-loop
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
