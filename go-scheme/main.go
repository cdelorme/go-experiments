package main

import (
	"fmt"
	"net/url"
)

func main() {
	u, _ := url.Parse("https://www.google.com")

	// @note: you can change the scheme using this property
	// u.Scheme = "http"

	fmt.Printf("Schema: %s\n", u.Scheme)
	fmt.Printf("Host: %s\n", u.Host)
	fmt.Printf("Whole: %+v\n", u)
}
