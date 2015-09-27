package main

import (
	"encoding/json"
	"fmt"

	"github.com/cdelorme/go-maps"
)

type User struct {
	Username string `json:"username"`
}

func main() {
	var data map[string]interface{}
	err := json.Unmarshal([]byte("{\"username\":\"Casey\"}"), &data)
	if err != nil {
		fmt.Printf("Failure!: %s\n", err)
		return
	}

	fmt.Printf("Data: %+v\n", data)

	user := &User{}
	maps.To(user, data)
	fmt.Printf("User: %+v\n", user)
}
