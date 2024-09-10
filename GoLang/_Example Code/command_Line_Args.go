// Interacting with the Operating Systems - implementation of Unix echo
package main

import (
	"fmt"
	"os" // os is a package for interacting with the OS
)

func main() {
	var s, sep string                   // variable declartions, initialised to "" (0 for int)
	for i := 1; i < len(os.Args); i++ { // standard for loop; no ( ); '{' in the same line
		s += sep + os.Args[i] // os.Args is a slide of strings (slide = sort of dynamic array), same as in Java
		sep = " "
	}
	fmt.Println(s)

	// use for-loop for most things -- a traditional while loop is written as
	// for condition {
	//	...
	// }
	// no condition ==> while true

}
