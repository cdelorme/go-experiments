package main

import (
	"flag"
	"log"
	"os"
)

var address = flag.String("address", ":10001", "Address of the server we are connecting to (defaults to localhost:10001)")

func main() {
	s := &Server{}
	if err := s.Init(*address); err != nil {
		log.Printf("error initializing: %s\n", err)
		os.Exit(1)
	}
	defer s.Close()
	log.Printf("%#v\n", s)
	s.Run()
}
