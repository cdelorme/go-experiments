package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var c = make(chan os.Signal)

func capture() {
	for _ = range c {
		fmt.Println("\nInterrupt Received!")
		os.Exit(0)
	}
}

func main() {
	signal.Notify(c, syscall.SIGHUP)
	go capture()

	fmt.Println("This will run until it receives a SIGHUP!")
	for {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf(".")
	}
}
