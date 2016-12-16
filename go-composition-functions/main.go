package main

// a quick demo of accessing inner functions
// and properties using composition

import (
	"fmt"
)

type inner struct {
	Password string
}

func (self *inner) do() {
	fmt.Println(self.Password)
}

type Outer struct {
	inner
}

func (self *Outer) Test() {
	self.do()
}

func main() {
	t := Outer{}
	t.Password = "dangajits"
	t.Test()
}
