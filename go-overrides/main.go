package main

// @note: (failed) experiment based on nodejs's string lib
// in node, S = require('strings');
// this allows similar casting via `S("a string")`
// the same can be done within go's type-strict system
// allowing us to append functionality with indirection
// but due to casting performance is (probably) abysmal
// hence this is marked as "failed"

import (
	"fmt"
	"strings"
)

type S string

func (s *S) Contains(substr string) bool {
	return strings.Contains(s.String(), substr)
}

func (s *S) String() string {
	return string(*s)
}

func main() {

	// this is how it usually looks in node
	var s = S("Bananas")
	// if it were a library, we'd probably have to do s.S() for package s

	// but this also works
	// var s S = "Bananas"

	// This is how we can use it
	fmt.Println(s.Contains("Ba"))
	fmt.Println(s.Contains("Fa"))

	// as opposed to:
	fmt.Println(strings.Contains("Bananas", "Ba"))
}
