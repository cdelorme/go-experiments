package main

import (
	"fmt"
)

// simple test case to verify parameters passed by reference vs value
// when using defer cannot be changed, thus if you want to set the
// value after calling defer you must use memory reference and
// the deferred operation must de-reference the value

func main() {

	// initialize but do not assign
	var data string

	// defer by value and by reference
	defer pv(data)
	defer pr(&data)

	// set value
	data = "banana"

	// the pv call will print an empty string
	// because it remembers the value at defer-time

	// the pr address is not changed, but what it
	// points to is, so by de-referencing in the
	// function, we can get the latest value
	// which let's us pass changing data to a defer
}

func pr(data *string) {
	fmt.Println("reference: ", *data)
}

func pv(data string) {
	fmt.Println("value: ", data)
}
