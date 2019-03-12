package main

import (
	"fmt"
)

func main() {

	// manual casting to and from bytes and uint32
	var num uint32 = 123
	bit := []byte{byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24)}
	var num2 uint32 = uint32(bit[0]) | uint32(bit[1])<<8 | uint32(bit[2])>>16 | uint32(bit[3])<<24

	fmt.Println("Bytes and Numbers")
	fmt.Printf("%d; %v; %d\n", num, bit, num2)

}
