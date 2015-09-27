package main

import (
	"fmt"
)

type T struct {
	Data string
}

func main() {

	// struct pointer slice (bad practice, pointer to pointer to data)
	// can only see pointer to T, not T itself
	ts := make([]*T, 0)
	ts = append(ts, &T{Data: "Some string value"})
	fmt.Printf("T's: %+v\n", ts)

	// map string interface /w struct
	// prints memory address, can only see key name
	// as well as type stored
	x := make(map[string]interface{}, 0)
	x["stuff"] = &ts
	x["more direct"] = &T{Data: "needs more strings!"}
	fmt.Printf("Things! %+v\n", x)
	fmt.Printf("Things! %#v\n", x)
	fmt.Printf("Things! %T\n", x)
}
