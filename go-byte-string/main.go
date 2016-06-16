package main

// test capturing a string as a byte array from json unmarshal
// @note: confirmed that while it may encode base64 correctly it does not seem to be readable from string inputs

import (
	"encoding/json"
	"fmt"
)

type Thing struct {
	Data []byte `json:"data,omitempty"`
}

func main() {
	toJson()
	fromJson()
}

func toJson() {
	t := Thing{Data: []byte("test")}
	d, e := json.Marshal(t)
	if e != nil {
		fmt.Printf("Error: %s\n", e)
	}
	fmt.Printf("%s\n", d)
}

func fromJson() {
	d := []byte(`{"data": "test"}`)
	t := Thing{}
	e := json.Unmarshal(d, &t)
	if e != nil {
		fmt.Printf("Error: %s\n", e)
	}
	fmt.Printf("%+v\n", t)

	// try base64 decode on t.Data?
	// fmt.Printf("%v\n", string(t.Data))
}
