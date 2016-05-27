package main

import (
	"fmt"
	"net"
)

// udp server demonstration
// accept and respond to data
// plans to add:
// - identity-pools by client address
// - test case to verify responding through NAT
// - message en/de-cryption

func main() {
	var err error

	serverAddr, err := net.ResolveUDPAddr("udp", ":10001")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	// set size constraint
	conn.SetReadBuffer(1024)
	conn.SetWriteBuffer(1024)

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		fmt.Printf("(%s): %s\n", addr, buf[0:n])
		go send(conn, addr)
	}
}

func send(c *net.UDPConn, addr *net.UDPAddr) {
	c.WriteToUDP([]byte("received"), addr)
}
