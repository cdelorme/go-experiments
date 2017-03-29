package chan_test

// an unscientific case to compare channels with different types
// to see whether one type outperforms another

import (
	"sync"
	"testing"
)

var concurrency = 4

func BenchmarkStructChan(b *testing.B) {
	ch := make(chan struct{})
	for i := 0; i < concurrency; i++ {
		go func(c chan struct{}) {
			for {
				<-c
			}
		}(ch)
	}
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
}

func BenchmarkBoolChan(b *testing.B) {
	ch := make(chan bool)
	for i := 0; i < concurrency; i++ {
		go func(c chan bool) {
			for {
				<-c
			}
		}(ch)
	}
	for i := 0; i < b.N; i++ {
		ch <- true
	}
}

func BenchmarkIntChan(b *testing.B) {
	ch := make(chan int)
	for i := 0; i < concurrency; i++ {
		go func(c chan int) {
			for {
				<-c
			}
		}(ch)
	}
	for i := 0; i < b.N; i++ {
		ch <- 1
	}
}
