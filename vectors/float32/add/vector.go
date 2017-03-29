package main

import (
	"math"
)

func main() {}

type Vector struct {
	X float32
	Y float32
}

func (self *Vector) Add(v *Vector) *Vector {
	return &Vector{X: self.X + v.X, Y: self.Y + v.Y}
}

func (self *Vector) Subtract(v *Vector) *Vector {
	return &Vector{X: self.X - v.X, Y: self.Y - v.Y}
}

func (self *Vector) Multiply(scalar float32) *Vector {
	return &Vector{X: self.X * scalar, Y: self.Y * scalar}
}

func (self *Vector) Divide(scalar float32) *Vector {
	v := &Vector{}
	if scalar > 0 {
		v.X = self.X / scalar
		v.Y = self.Y / scalar
	}
	return v
}

func (self *Vector) Length() float32 {
	return float32(math.Sqrt(float64(self.X*self.X) + float64(self.Y*self.Y)))
}
