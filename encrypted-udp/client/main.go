package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

var address = flag.String("address", "127.0.0.1:10001", "Address ofn the server we are connecting to")
var identity = flag.String("username", "", "Name to show in chat")

func main() {
	flag.Parse()

	c := &Client{}
	if err := c.Init(*identity, *address); err != nil {
		log.Printf("error initializing: %s\n", err)
		os.Exit(1)
	}
	defer c.Close()
	log.Printf("%#v\n", c)

	go c.Receive()

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "quit" || message == "exit" {
			log.Printf("exiting...\n")
			break
		}
		if err := c.MessageSend(message); err != nil {
			log.Printf("error sending: %s\n", err)
		}
	}
}
