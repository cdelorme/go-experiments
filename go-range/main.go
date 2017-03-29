package main

import "fmt"

func main() {

	test := []string{"test", "me", "again", "please"}

	// @note: range sets the index per iteration complicating manipulation
	fmt.Println("with range:")
	for i, a := range test {
		if a == "me" {
			i += 2
			continue
		}
		fmt.Printf("Index: %d, Value: %s\n", i, a)
	}

	// @note: index manipulation is valid
	fmt.Println("\nwith index:")
	for i := 0; i < len(test); i++ {
		if test[i] == "me" {
			i++
			continue
		}
		fmt.Printf("Index: %d, Value: %s\n", i, test[i])
	}
}
