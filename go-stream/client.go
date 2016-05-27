package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

// udp client connection scenario
// plans to add:
// - latency computation
// - identifier per client
// - (standard?) ssl encryption

func main() {
	var err error

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	defer conn.Close()

	// we can ensure buffer size to ensure no data ends up getting lost
	conn.SetReadBuffer(1024)
	conn.SetWriteBuffer(1024)

	// async sends with sync receiver (to prevent main thread termination)
	go sender(conn)
	receiver(conn)
}

func sender(c *net.UDPConn) {
	// @todo: generate identifier for testing multiple clients
	// @todo: send timestamp instead to calculate latency
	i := 0
	for {
		if _, err := c.Write([]byte(strconv.Itoa(i))); err != nil {
			fmt.Printf("%s\n", err)
		}
		time.Sleep(time.Second * 1)
		i++
	}
}

func receiver(c *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		if _, err := bufio.NewReader(c).Read(buf); err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		fmt.Printf("%s\n", buf)
	}

	// stop-gap for sleep forever
	// select {}
}
