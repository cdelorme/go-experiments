package main

// a simple case demonstrating using a single test case as a mock
// for using exec with expectations (in this case pass or fail)

import (
	"fmt"
	"os/exec"
)

var cmd = exec.Command

func main() {
	e := cmd("echo", "hello world!")
	out, err := e.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("Output: %s", out)
}
