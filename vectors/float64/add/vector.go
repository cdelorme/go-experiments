package main

import (
	"math"
)

func main() {}

type Vector struct {
	X float64
	Y float64
}

func (self *Vector) Add(v *Vector) *Vector {
	return &Vector{X: self.X + v.X, Y: self.Y + v.Y}
}

func (self *Vector) Subtract(v *Vector) *Vector {
	return &Vector{X: self.X - v.X, Y: self.Y - v.Y}
}

func (self *Vector) Multiply(scalar float64) *Vector {
	return &Vector{X: self.X * scalar, Y: self.Y * scalar}
}

func (self *Vector) Divide(scalar float64) *Vector {
	v := &Vector{}
	if scalar > 0 {
		v.X = self.X / scalar
		v.Y = self.Y / scalar
	}
	return v
}

func (self *Vector) Length() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y)
}
