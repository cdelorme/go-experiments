package main

import (
	"errors"
	"fmt"
)

// can use package-level var (seen in std lib's aka idiomatic go?)
var MyError = errors.New("A package-level error.")

// const does not work, only primitives for const values
// const MyError = errors.New("A constant error.")

func main() {
	fmt.Printf("%s\n", MyError)
}
