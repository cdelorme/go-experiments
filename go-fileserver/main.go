package main

import (
	"net/http"
)

// simple example of how to serve files from go in 2 lines of code

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public/")))
	http.ListenAndServe(":8123", nil)
}
