package main

import "net"

// The client from the perspective of the server.
type Client struct {
	a        *net.UDPAddr
	key      [32]byte
	identity string
}
