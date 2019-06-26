package main

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/rand"
// 	"io"
// 	"testing"
// )

// func TestECDSA(t *testing.T) {
// 	message := make([]byte, ECDSAMessageSize)
// 	if _, err := io.ReadFull(rand.Reader, message); err != nil {
// 		t.Fatalf("failed to create message: %s", err)
// 	}

// 	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		t.Fatalf("failed to create rsa key: %s", err)
// 	}

// 	// ciphertext, err := RSAOAEPEncrypt(&key.PublicKey, message)
// 	// if err != nil {
// 	// 	t.Fatalf("failed to rsa oaep encrypt: %s", err)
// 	// } else if len(ciphertext) != 256 {
// 	// 	t.Fatal("failed to produce expected ciphertext length...")
// 	// }

// 	// data, err := RSAOAEPDecrypt(key, ciphertext)
// 	// if err != nil {
// 	// 	t.Fatalf("failed to rsa oaep decrypt: %s", err)
// 	// } else if !bytes.Equal(data, message) {
// 	// 	t.Fatal("failed to produce valid decrypted bytes")
// 	// }
// }
