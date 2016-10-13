package main

import (
	"math"
)

func main() {}

type Vector struct {
	X float32
	Y float32
}

func (self *Vector) Add(v *Vector) {
	self.X += v.X
	self.Y += v.Y
}

func (self *Vector) Subtract(v *Vector) {
	self.X -= v.X
	self.Y -= v.Y
}

func (self *Vector) Multiply(scalar float32) {
	self.X *= scalar
	self.Y *= scalar
}

func (self *Vector) Divide(scalar float32) {
	if scalar != 0 {
		self.X /= scalar
		self.Y /= scalar
	}
}

func (self *Vector) Length() float32 {
	return float32(math.Sqrt(float64(self.X*self.X) + float64(self.Y*self.Y)))
}
