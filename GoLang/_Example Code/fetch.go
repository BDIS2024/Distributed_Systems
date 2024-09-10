// Fetch prints the content found at a URL.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] { // looks at all url addresses provided in the command line
		fmt.Println("looking at address " + url)
		resp, err := http.Get(url) // using net/http package, it returns a pair -- first element is a struct
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body) // resp.Body is a field containing the response as a readable stream
		// ioutil.ReadAll reads the entire response (entire stream)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}

		fmt.Printf("%s", b)
	}
}
