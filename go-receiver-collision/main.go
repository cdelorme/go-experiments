// A quick verification that receivers do not account for pointer or
// non-pointer and will collide at the function name.
package main

import (
	"fmt"
)

type Vector struct {
	X float32
	Y float32
}

func (self *Vector) Add(v *Vector) {
	self.X += v.X
	self.Y += v.Y
}

func (self Vector) Add(v *Vector) Vector {
	self.X += v.X
	self.Y += v.Y
	return self
}

func main() {
	v := Vector{}
	fmt.Printf("%+v\n", v)
}
