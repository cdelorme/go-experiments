package main

import "fmt"

func main() {

	// create a string slice, and cast to interface
	data := make([]string, 0)
	data = append(data, "banana")
	data = append(data, "hammock")
	data = append(data, "sand")
	data = append(data, "witch")
	fmt.Printf("%+v\n", data)

	// cast to interface
	var dataTwo interface{}
	dataTwo = data
	fmt.Printf("%+v\n", dataTwo)
	// data appears to be the same

	// recast back to expected type
	dataThree := dataTwo.([]string)
	fmt.Printf("%+v\n", dataThree)
	// looks unchanged
}
