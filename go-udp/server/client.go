package main

import (
	"net"
)

// This is the server-side representation of a client
type client struct {
	a *net.UDPAddr
}
