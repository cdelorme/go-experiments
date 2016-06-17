package main

// verifies that a UDPAddr resolves to an IP and not a raw DNS string
// for 1:1 connections where public server ip is directly mapped that's great
// however it may not enough if the server sits behind a load balancer

import (
	"fmt"
	"net"
)

func main() {

	// requires a port
	dns := "www.google.com:80"

	// can catch errors, like missing a port
	addr, e := net.ResolveUDPAddr("udp", dns)
	if e != nil {
		fmt.Printf("Something went horribly wrong: %s\n", e)
		return
	}

	// somewhat unsurprisingly does resolve an IP
	fmt.Printf("New Address: %s (%+v)\n", addr, addr)
}
