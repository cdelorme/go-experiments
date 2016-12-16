package main

import "fmt"

type dumb struct {
	Property string
}

type dumber struct {
	Property string
}

// @note: even if you change the order of properties, composite is dealt with first
type smart struct {
	dumb
	Named dumber
}

func main() {
	s := smart{}
	s.Property = "composition takes priority"
	fmt.Printf("%#v\n", s)
}
