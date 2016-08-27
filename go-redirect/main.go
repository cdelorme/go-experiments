package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	url := "https" + strings.TrimPrefix(req.URL.String(), "http")
	fmt.Println(url)
}
