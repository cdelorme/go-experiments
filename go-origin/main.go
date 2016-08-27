package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Origin: %s\n", r.Header.Get("Origin"))
		fmt.Printf("Referer: %s\n", r.Referer())
	})
	fmt.Printf("Broke Down: %s", http.ListenAndServe(":3000", nil))
}
