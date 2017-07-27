// A microbenchmark to demonstrate the cost of pointers anonymous functions.
package main

import (
	"testing"
)

type Data struct {
	Value int
	value int
}

func (self *Data) GetValue() int {
	return self.value
}

type DataTwo struct {
	Value *int
	value *int
}

func (self *DataTwo) GetValue() *int {
	return self.value
}

type DataThree struct{}

func (self *DataThree) Do()           {}
func (self *DataThree) DoParam(_ int) {}

func Do()           {}
func DoParam(_ int) {}

var thing = Data{Value: 42, value: 42}
var thingTwo = DataTwo{Value: &thing.Value, value: &thing.Value}
var thingThree = DataThree{}
var temp int
var tempTwo *int

func BenchmarkProperty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		temp = thing.Value
	}
}

func BenchmarkAccessor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		temp = thing.GetValue()
	}
}

func BenchmarkPropertyPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tempTwo = thingTwo.Value
	}
}

func BenchmarkAccessorPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tempTwo = thingTwo.GetValue()
	}
}

func BenchmarkMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		thingThree.Do()
	}
}

func BenchmarkMethodParam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		thingThree.DoParam(i)
	}
}

func BenchmarkFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Do()
	}
}

func BenchmarkFuncParam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoParam(i)
	}
}

func BenchmarkAnonFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {}()
	}
}

func BenchmarkAnonFuncParam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func(_ int) {}(i)
	}
}
