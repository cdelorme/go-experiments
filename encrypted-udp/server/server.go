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

var errInvalidKeySize = fmt.Errorf("keys must be %d bytes...", KeySize)

// A server implementation that creates a goroutine per
// connection and passes a shared channel to communicate.
type Server struct {
	priv    [32]byte
	pub     [32]byte
	c       *net.UDPConn
	clients map[string]Client
}

// Clear all clients and close the server.
func (s *Server) Close() error {
	s.clients = make(map[string]Client, 0)
	return s.c.Close()
}

// Parse the address and start the server with chosen buffer size.
func (s *Server) Init(address string) error {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	copy(s.priv[:], priv[:])
	copy(s.pub[:], pub[:])

	if serverAddr, err := net.ResolveUDPAddr("udp", address); err != nil {
		return err
	} else if s.c, err = net.ListenUDP("udp", serverAddr); err != nil {
		return err
	}

	s.c.SetReadBuffer(BufferSize)
	s.c.SetWriteBuffer(BufferSize)
	s.clients = make(map[string]Client, 0)
	return nil
}

// Listen for new connections to create clients with their own goroutines.
//
// Since UDP is "connectionless" we may want to track the last received
// message time so we can clear "idle" clients after a fixed time.
//
// Writing a resilient client to handle disconnection would be interesting.
func (s *Server) Run() {
	b := make([]byte, BufferSize)
	for {
		l, addr, err := s.c.ReadFromUDP(b)
		if err != nil {
			log.Printf("failed to read from connection: %s\n", err)
			continue
		}
		go s.MessageProcess(addr, b[:l])
	}
}

func (s *Server) MessageProcess(addr *net.UDPAddr, message []byte) {
	if len(message) < len(Signature)+1 || !bytes.Equal(message[:len(Signature)], Signature[:]) {
		log.Printf("Signature does not match for address %s, discarding...\n", addr.String())
		return
	}
	switch message[len(Signature)] {
	case MessageHandshake:
		s.HandshakeReceive(addr, message[len(Signature)+1:])
	case MessageChat:
		s.MessageReceive(addr, message[len(Signature)+1:])
	default:
		log.Printf("unknown message type (%d) from address: %s\n", message[len(Signature)], addr.String())
	}
}

func (s *Server) HandshakeReceive(addr *net.UDPAddr, shake []byte) {
	if len(shake) < KeySize {
		log.Printf("handshake failed due to key size (%d): %s", len(shake), errInvalidKeySize)
		s.Disconnected(addr, "invalid key size...")
		return
	} else if len(shake) > KeySize+MaxIdentitySize {
		log.Printf("identity too large (%d: %s), sending disconnected...\n", len(shake[KeySize:]), string(shake[KeySize:]))
		s.Disconnected(addr, "identity too large...")
		return
	}

	var pub [KeySize]byte
	copy(pub[:], shake[:KeySize])

	identity := make([]byte, len(shake)-KeySize)
	copy(identity, shake[KeySize:len(shake)])

	c := Client{a: addr, identity: string(identity)}
	box.Precompute(&c.key, &pub, &s.priv)

	data := append(append(Signature[:], MessageHandshake), s.pub[:]...)
	if _, err := s.c.WriteToUDP(data, addr); err != nil {
		log.Printf("failed to write handshake message to %s: %s", addr.String(), err)
		s.Disconnected(addr, "failed to send handshake reply...")
		return
	}

	s.clients[addr.String()] = c
	log.Printf("Established connection with %#v\n", c)
}

// Sends a disconnected message to the client, with a reason that might
// be useful with an interactive interface or for debugging.
func (s *Server) Disconnected(addr *net.UDPAddr, reason string) {
	data := append(append(Signature[:], MessageDisconnected), []byte(reason)...)
	s.c.WriteToUDP(data, addr)
}

func (s *Server) MessageReceive(addr *net.UDPAddr, ciphertext []byte) {
	if _, ok := s.clients[addr.String()]; !ok {
		log.Printf("No registered client %s\n", addr.String())
		s.Disconnected(addr, "not registered...")
		return
	}

	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])
	key := s.clients[addr.String()].key
	message, ok := box.OpenAfterPrecomputation(nil, ciphertext[24:], &nonce, &key)
	if !ok {
		log.Printf("failed to decrypt chat message from %s\n", addr.String())
		s.Disconnected(addr, "failed to decrypt chat message...")
		return
	}

	if len(message) > MaxMessageSize {
		log.Printf("message received from %s is too large: %s\n", addr.String, string(message))
	}

	if _, err := rand.Read(nonce[:]); err != nil {
		log.Printf("failed to generate new nonce: %s\n", err)
		return
	}

	message = append([]byte(s.clients[addr.String()].identity+": "), message...)

	// @note: ideally this would send each message on a goroutine,
	// but simply prefixing with go allows the client to change in the loop
	// which breaks.
	//
	// @note: even with goroutines there may be a bias as to who receives
	// first due to the map order.
	//
	// @note: this may not be concurrently safe since new connections can occur
	// in parallel to sending messages, which may lead to a race condition on
	// the array of clients.
	log.Printf("Sending %s to clients: %#v", message, s.clients)
	for _, client := range s.clients {
		s.MessageSend(&client, nonce, message)
	}
}

func (s *Server) MessageSend(c *Client, nonce [24]byte, message []byte) {
	ciphertext := box.SealAfterPrecomputation(nonce[:], message, &nonce, &c.key)
	log.Printf("Identity: %s", c.identity)
	data := append(append(Signature[:], MessageChat), ciphertext...)
	s.c.WriteToUDP(data, c.a)
}
