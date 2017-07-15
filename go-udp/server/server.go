package main

import (
	"log"
	"net"
)

const bufferSize = 1024

// A server implementation with the ability to track multiple clients.
//
// In the future it may make sense to benchmark slice search vs maps
// with up to 10,000 elements.
type server struct {
	c       *net.UDPConn
	clients map[string]client
}

// Establish a connection on the supplied address, and set
// expected buffer sizes to avoid unpredictable behavior.
//
// Equal size buffers avoids DoS concerns.
//
// Finally initializes the map of clients which will be
// used to distribute messages.
func (s *server) Init(address string) error {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	s.c, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	s.c.SetReadBuffer(bufferSize)
	s.c.SetWriteBuffer(bufferSize)
	s.clients = make(map[string]client, 0)
	return nil
}

// Abstraction to close the established connection.
func (s *server) Close() {
	s.c.Close()
}

// Prepares a reusable buffer to reduce allocations, and infinitely loops
// reading from the connection.
//
// It collects client representations by address in string format, so
// it can distribute received messages all all clients.
//
// The distribution method leverages goroutines so each send is not
// waiting on the other receive.
//
// Keep in mind this has no protection from garbage input.
//
// The distribution method, while concurrent, will be biased by the
// earliest element in the map when iterated.
func (s *server) Run() {
	b := make([]byte, bufferSize)
	for {
		n, addr, err := s.c.ReadFromUDP(b)
		if err != nil {
			log.Printf("error receiving: %s\n", err)
			continue
		} else if n == 0 {
			continue
		}
		if _, ok := s.clients[addr.String()]; !ok {
			s.clients[addr.String()] = client{a: addr}
		}
		for i := range s.clients {
			if s.clients[i].a.String() == addr.String() {
				continue
			}
			go s.Send(b[:n], s.clients[i])
		}
	}
}

// Accepts the message in byte array format and a client to
// send the message to.
func (s *server) Send(message []byte, c client) {
	n, err := s.c.WriteToUDP(message, c.a)
	if err != nil {
		log.Printf("failed to send: %s\n", err)
	} else if n != len(message) {
		log.Printf("expected to write %d bytes, but wrote %d instead...", len(message), n)
	}
}
