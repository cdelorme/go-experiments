package main

import (
	"fmt"
)

func main() {
	// to verify that functions called in the return line happen before deferred statements
	// also that defer is a lifo stack
	TestDeferReturnOrder()
}

func TestDeferReturnOrder() (int, error) {
	defer fmt.Println("Defer 1")
	defer fmt.Println("Defer 2")
	defer fmt.Println("Defer 3")
	defer fmt.Println("Defer 4")
	return fmt.Println("Return")
}
