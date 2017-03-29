package main

import (
	"fmt"
)

func main() {
	args := []interface{}{"one", struct{ Two int }{Two: 2}, 3}

	tmp(args...)

	// @note: you can manipulate a slice before passing
	// tmp(args[0:]...)
}

func tmp(args ...interface{}) {
	fmt.Printf("Args: %+v", args)
}
