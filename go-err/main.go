package main

import (
	"errors"
	"fmt"
)

func main() {
	err := errors.New("Testing error strings")
	fmt.Printf("Auto-converted: %s\n", err)
	fmt.Printf("Forced conversion: %s\n", err.Error())
}
