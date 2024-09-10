package main

import "fmt" // pronounced “fumpt” - formatting package

func main() {
	fmt.Println("Hello world!") // similar System.out.println in Java
	fmt.Println("Hello world!")
	fmt.Println("Hello world!")
}

/*
	Things to remember:
	- main is a special keyword (like in Java)
	- Go has packages (like libraries), 'import' is used to import other packages
	- ';' is not as important as in Java, only for separating commands on the same line
	- the opening '{' must be on the same line as the end of the func declaration
	- 'gofmt' tool helps reformatting -- in VS Code, the Go extension does it for you on save
*/
