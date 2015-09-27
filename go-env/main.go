package main

import (
	"fmt"
	"os"
)

func main() {

	// can access/print all environment variables
	envs := os.Environ()
	for _, v := range envs {
		fmt.Printf("%s\n", v)
	}
}
