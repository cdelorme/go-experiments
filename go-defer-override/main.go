package main

// a quick verification to demonstrate that defer parameters are accepted
// verbatim, and to override requires a layer that crosses over scope

import "fmt"

func main() {
	message := "Hello World"

	defer fmt.Println(message)
	defer func() {
		fmt.Println(message)
	}()

	message = "Goodbye World"
}
