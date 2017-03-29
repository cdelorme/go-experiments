package main

import (
	"math"
)

func main() {}

type Vector struct {
	X float64
	Y float64
}

func (self *Vector) Add(v *Vector) {
	self.X += v.X
	self.Y += v.Y
}

func (self *Vector) Subtract(v *Vector) {
	self.X -= v.X
	self.Y -= v.Y
}

func (self *Vector) Multiply(scalar float64) {
	self.X *= scalar
	self.Y *= scalar
}

func (self *Vector) Divide(scalar float64) {
	if scalar != 0 {
		self.X /= scalar
		self.Y /= scalar
	}
}

func (self *Vector) Length() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y)
}
