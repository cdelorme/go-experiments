package main

import (
	"testing"
)

func BenchmarkVector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := &Vector{X: 10, Y: 12}
		v = v.Add(&Vector{X: 12, Y: 10})
		v = v.Subtract(&Vector{X: 12, Y: 10})
		v = v.Multiply(3)
		v = v.Divide(3)
		v.Length()
	}
}
