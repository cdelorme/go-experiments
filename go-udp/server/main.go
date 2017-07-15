package main

import (
	"flag"
	"log"
)

var address = flag.String("address", ":10001", "Address ofn the server we are connecting to")

func main() {
	s := &server{}
	if err := s.Init(*address); err != nil {
		log.Printf("error initializing: %s\n", err)
		return
	}
	defer s.Close()
	log.Printf("%#v\n", s)
	s.Run()
}
