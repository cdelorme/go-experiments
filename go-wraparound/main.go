package main

// test case demonstrates simple min/max computation
// and tests wrap effect when exceeding the size

import "fmt"

// compute min and max values for signed and unsigned integers
const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func main() {

	// print the min and max values
	fmt.Printf("Max Unsigned Int: %d\n", MaxUint)
	fmt.Printf("Min Unsigned Int: %d\n", MinUint)

	fmt.Printf("Max Int: %d\n", MaxInt)
	fmt.Printf("Min Int: %d\n", MinInt)

	// test unsigned integer wrap effect
	var ui uint
	fmt.Printf("Default Unsigned Integer is Minimum: %d\n", ui)

	ui--
	fmt.Printf("Minimum Negative Wrap is Maximum: %d\n", ui)

	ui++
	fmt.Printf("Maximum Positive Wrap is Minimum: %d\n", ui)

	// test signed integer wrap effect
	var i int
	fmt.Printf("Default Signed Integer is middle: %d\n", i)

	i = MinInt
	i--
	fmt.Printf("Minimum Signed Integer Negative Wrap is Maximum: %d\n", i)

	i++
	fmt.Printf("Maximum Signed Integer Negative Wrap is Minimum: %d\n", i)
}
