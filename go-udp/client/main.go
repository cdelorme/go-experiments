package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var address = flag.String("address", "127.0.0.1:10001", "Address ofn the server we are connecting to")
var identity = flag.String("username", "", "Name to show in chat")

func main() {
	c := &client{}
	if err := c.Init(*identity, *address); err != nil {
		log.Printf("error initializing: %s\n", err)
		return
	}
	defer c.Close()
	log.Printf("%#v\n", c)

	go func(c *client) {
		for {
			d, err := c.Receive()
			if err != nil {
				log.Printf("error receiving: %s\n", err)
				continue
			}
			fmt.Println(d)
		}
	}(c)

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "quit" || message == "exit" {
			log.Printf("exiting...\n")
			break
		}
		if err := c.Send(message); err != nil {
			log.Printf("error sending: %s\n", err)
		}
	}
}
