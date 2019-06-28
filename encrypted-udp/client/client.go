package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/nacl/box"
)

// Shared constants representing a known message types,
// sizes used to establish buffers, and service signature
// to easily drop unknown traffic.
const (
	MessageHandshake byte = iota
	MessageDisconnected
	MessageChat
)

const (
	BufferSize      = 508
	NaClPadding     = 16
	NaClNonceSize   = 24
	MaxIdentitySize = 20
	MaxMessageSize  = BufferSize - (NaClNonceSize + NaClPadding + MaxIdentitySize + 2) // 2 spaces for formatting
	KeySize         = 32
)

var Signature = [...]byte{1, 2, 3, 4}

var errNoIdentity = fmt.Errorf("Identity is empty!")
var errMessageTooBig = fmt.Errorf("messages must be under %d characters...", BufferSize)
var errIdentityTooBig = fmt.Errorf("username must be under %d bytes...", MaxIdentitySize)
var errKeyTooBig = fmt.Errorf("keys must be %d bytes...", KeySize)

type Client struct {
	identity string
	c        *net.UDPConn
	key      [32]byte
}

// Reads from the connection, checking the signature, and using the type
// to decide where to send the content.
//
// Errors will be logged.
func (c *Client) Receive() {
	b := make([]byte, BufferSize)

	for {
		l, err := c.c.Read(b)
		if err != nil {
			log.Printf("failed to read from connection: %s\n", err)
			continue
		}
		go c.MessageProcess(b[:l])
	}
}

func (c *Client) MessageProcess(message []byte) {
	if len(message) < len(Signature)+1 || !bytes.Equal(message[:len(Signature)], Signature[:]) {
		log.Printf("Signature does not match, discarding...\n")
		return
	}
	switch message[len(Signature)] {
	case MessageDisconnected:
		c.HandshakeSend()
	case MessageHandshake:
		c.HandshakeReceive(message[len(Signature)+1:])
	case MessageChat:
		c.MessageReceive(message[len(Signature)+1:])
	default:
		log.Printf("unknown message type: %d\n", message[len(Signature)])
	}
}

// Close the UDP connection.
func (c *Client) Close() {
	if c.c != nil {
		c.c.Close()
	}
}

// Establishes the UDP connection to a specified remote server
// using a defined buffer size.
//
// Establishes identity, and sets up asynchronous listener.
//
// Sends handshake request to establish the connection.
func (c *Client) Init(identity, address string) error {
	c.identity = identity
	if c.identity == "" {
		return errNoIdentity
	} else if len([]byte(c.identity)) > MaxIdentitySize {
		return errIdentityTooBig
	}

	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return err
	}

	c.c, err = net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return err
	}
	c.c.SetReadBuffer(BufferSize)
	c.c.SetWriteBuffer(BufferSize)

	return c.HandshakeSend()
}

// Used when establishing a connection, or dealing with
// disconnection.
func (c *Client) HandshakeSend() error {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// @note: adding RSA encryption would prevent a proxy MITM attacker
	// from capturing and modifying all traffic in transit, or worse yet
	// setting your identity to mrpoopybutthole.
	//
	// Signature validation would not necessarily be required for processing
	// the received handshake, since an attacker could only fudge the response
	// data breaking your connection.

	// copy private key to precompute when we get the return handshake
	copy(c.key[:], priv[:])

	// prepare a message the the public key and identity
	data := append(append(append(Signature[:], MessageHandshake), pub[:]...), []byte(c.identity)...)

	_, err = c.c.Write(data)
	return err
}

// Complete the handshake by precomputing the received key.
func (c *Client) HandshakeReceive(key []byte) {
	if len(key) != KeySize {
		log.Printf("handshake failed due to key size (%d): %s", len(key), errKeyTooBig)
		return
	}

	var priv [32]byte
	copy(priv[:], c.key[:])

	var pub [32]byte
	copy(pub[:], key)

	box.Precompute(&c.key, &pub, &priv)
	log.Printf("Handshake completed!\n")
}

func (c *Client) MessageReceive(ciphertext []byte) {
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])
	message, ok := box.OpenAfterPrecomputation(nil, ciphertext[24:], &nonce, &c.key)
	if !ok {
		log.Printf("failed to decrypt...\n")
		return
	}
	fmt.Println(string(message))
}

func (c *Client) MessageSend(message string) error {
	// @note: it may be sensible to add a flag that is set during HandshakeSend
	// which can be checked here if the handshake is incomplete and resend it.
	//
	// Since UDP is "connectionless", unless we receive an explicit command for
	// disconnection we won't try to establish a new handshake, but the client
	// is written to be resilient so when the server cannot decrypt it will
	// send an appropriate response to trigger the reconnection, at the cost of
	// lost inbound messages sent by the client.

	messageBytes := []byte(message)
	if len(messageBytes) > MaxMessageSize {
		return errMessageTooBig
	}

	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return err
	}

	ciphertext := box.SealAfterPrecomputation(nonce[:], messageBytes, &nonce, &c.key)

	data := append(append(Signature[:], MessageChat), ciphertext...)

	_, err := c.c.Write(data)
	return err
}
