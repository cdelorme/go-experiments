package main

// verifying making a thread wait for a channel response and close

import (
	"fmt"
	"time"
)

type internal struct {
	Value int
}

type test struct {
	Test *internal
}

func main() {

	c := make(chan test)

	go func() {
		for i := range c {
			time.Sleep(1 * time.Second)
			fmt.Printf("%d\n", i.Test.Value)
		}
	}()
	fmt.Println("Waiting...")

	i := &internal{Value: 1}
	v := test{Test: i}
	c <- v
	i.Value = 5

	time.Sleep(5 * time.Second)
	close(c)
	fmt.Println("Done...")
}
