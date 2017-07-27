// A demonstration that you can bind to custom host names by modifying
// your systems hosts file (eg. `/etc/hosts`).
//
// Try try it add "127.0.0.1 www.example.com" to your hosts file.
//
// Some situations where this can be very useful include predictable
// addresses for private service communication, as well as mocking
// dependencies.
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love serving %s!", r.Host)
}

func main() {
	http.HandleFunc("/", handler)
	e := http.ListenAndServe("www.example.com:8080", nil)
	fmt.Printf("%s\n", e)
}
