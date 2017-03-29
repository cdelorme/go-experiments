package main

import (
	"fmt"
	"time"
)

type Data struct {
	Interval int
	cb       func()
}

func main() {
	fmt.Printf("This will end after 10 seconds...")
	d := Data{}

	d.cb = func() {
		for {
			time.Sleep(200 * time.Millisecond)

			// @note: this is a race condition waiting to happen
			fmt.Printf("Interval: %d\n", d.Interval)
		}
	}
	go d.cb()

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		d.Interval++
	}
	time.Sleep(1 * time.Second)
}
