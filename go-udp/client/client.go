package main

// A rudimentary client implementation
//
// Note that read and write on a connection without a
// timeout (deadline) will block indefinetally.

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const bufferSize = 1024

var errTooBig = errors.New("messages must be under 1024 characters...")

type client struct {
	identity string
	c        *net.UDPConn
	w        io.Writer
	r        io.Reader
}

// Translates the supplied address to a UDP format, and
// acquires a free local UDP address.
//
// Sets the identity if empty to a GUID.
//
// Establishes a connection to the server, and restricts
// buffer size for predictable behavior.
func (c *client) Init(identity, address string) error {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return err
	}

	c.identity = identity
	if c.identity == "" {
		c.identity = GUID()
	}

	c.c, err = net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return err
	}

	c.c.SetReadBuffer(bufferSize)
	c.c.SetWriteBuffer(bufferSize)

	// @todo: establish deadlines to timeout read and write
	// only necessary if we want to terminate when the server
	// is unavailable or becomes unavailable
	// c.c.SetReadDeadline(time.Now().Add(timeoutDuration))
	// c.c.SetWriteDeadline(time.Now().Add(timeoutDuration))

	return nil
}

// Abstraction to close the established connection.
func (c *client) Close() {
	c.c.Close()
}

// Accepts a message in string format, appends the identity, and
// sends it over the UDP connection, verifying no errors and the
// correct number of bytes were written.
func (c *client) Send(message string) error {
	send := c.identity + ": " + message
	if len(send) > bufferSize {
		return errTooBig
	}
	n, err := c.c.Write([]byte(send))
	if n != len(send) {
		return fmt.Errorf("Expected to send %d bytes, but send %d instead", len(send), n)
	}
	return err
}

// Reads from the connection, passing back both the results and any
// errors encountered.
func (c *client) Receive() (string, error) {
	b := make([]byte, bufferSize)
	l, err := c.c.Read(b)
	return string(b[:l]), err
}
