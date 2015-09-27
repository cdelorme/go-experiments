package main

import "fmt"

type D struct {
	Name string
}

func main() {

	// create a struct
	s := D{Name: "Casey"}

	// create an unsigned pointer int to our struct
	var p uintptr = &s

	// what does that look like?
	fmt.Printf("D: %+v\n", p)

	// conclusion: this doesn't work, requires "unsafe" package
	// casting to unsafe.Pointer?  Or calling unsafe.Pointer()?
	// seems to deal with C pointers, probably useful for cgo
}
