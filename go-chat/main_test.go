package main

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestPlacebo(_ *testing.T) {}

func BenchmarkChannelOne(b *testing.B) {
	c := ChannelChat{}
	c.Add(&Connection{stream: ioutil.Discard})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}

func BenchmarkMutexOne(b *testing.B) {
	c := Chat{}
	c.Add(&ConnectionMutex{Connection: Connection{stream: ioutil.Discard}})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}

func BenchmarkChannelTwenty(b *testing.B) {
	c := ChannelChat{}
	for i := 0; i < 20; i++ {
		c.Add(&Connection{stream: ioutil.Discard})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}

func BenchmarkMutexTwenty(b *testing.B) {
	c := Chat{}
	for i := 0; i < 20; i++ {
		c.Add(&ConnectionMutex{Connection: Connection{stream: ioutil.Discard}})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}

func BenchmarkChannelHundred(b *testing.B) {
	c := ChannelChat{}
	for i := 0; i < 100; i++ {
		c.Add(&Connection{stream: ioutil.Discard})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}

func BenchmarkMutexHundred(b *testing.B) {
	c := Chat{}
	for i := 0; i < 100; i++ {
		c.Add(&ConnectionMutex{Connection: Connection{stream: ioutil.Discard}})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}

func BenchmarkChannelThousand(b *testing.B) {
	c := ChannelChat{}
	for i := 0; i < 1000; i++ {
		c.Add(&Connection{stream: ioutil.Discard})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}

func BenchmarkMutexThousand(b *testing.B) {
	c := Chat{}
	for i := 0; i < 1000; i++ {
		c.Add(&ConnectionMutex{Connection: Connection{stream: ioutil.Discard}})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Send(Message{
			Content:   "This is a demonstration",
			Author:    "Casey",
			Timestamp: time.Now(),
		})
	}
}
