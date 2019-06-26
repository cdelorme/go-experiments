package main

// A way to test rsa key pair encryption and decryption.
//
// Also benchmarks to verify performance of each phase.

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"testing"
)

func TestRSA(t *testing.T) {
	message := make([]byte, RSAMessageSize)
	if _, err := io.ReadFull(rand.Reader, message); err != nil {
		t.Fatalf("failed to create message: %s", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		t.Fatalf("failed to create rsa key: %s", err)
	}

	ciphertext, err := RSAOAEPEncrypt(&key.PublicKey, message)
	if err != nil {
		t.Fatalf("failed to rsa oaep encrypt: %s", err)
	} else if len(ciphertext) != 256 {
		t.Fatal("failed to produce expected ciphertext length...")
	}

	data, err := RSAOAEPDecrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("failed to rsa oaep decrypt: %s", err)
	} else if !bytes.Equal(data, message) {
		t.Fatal("failed to produce valid decrypted bytes")
	}
}

// Benchmark for RSA OAEP encryption
func BenchmarkRSAEncrypt(b *testing.B) {
	message := make([]byte, RSAMessageSize)
	if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		b.Fatalf("failed to create rsa key: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ciphertext, err := RSAOAEPEncrypt(&key.PublicKey, message)
		if err != nil {
			b.Fatalf("failed to rsa oaep encrypt: %s", err)
		} else if len(ciphertext) != 256 {
			b.Fatal("failed to produce expected ciphertext length...")
		}
	}
}

// Benchmark for RSA OAEP decryption
func BenchmarkRSADecrypt(b *testing.B) {
	message := make([]byte, RSAMessageSize)
	if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		b.Fatalf("failed to create rsa key: %s", err)
	}

	ciphertext, err := RSAOAEPEncrypt(&key.PublicKey, message)
	if err != nil {
		b.Fatalf("failed to rsa oaep encrypt: %s", err)
	} else if len(ciphertext) != 256 {
		b.Fatal("failed to produce expected ciphertext length...")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := RSAOAEPDecrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("failed to rsa oaep decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("failed to produce valid decrypted bytes")
		}
	}
}

// Benchmark for RSA OAEP encryption followed by decryption
func BenchmarkRSA(b *testing.B) {
	message := make([]byte, RSAMessageSize)
	if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		b.Fatalf("failed to create rsa key: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ciphertext, err := RSAOAEPEncrypt(&key.PublicKey, message)
		if err != nil {
			b.Fatalf("failed to rsa oaep encrypt: %s", err)
		} else if len(ciphertext) != 256 {
			b.Fatal("failed to produce expected ciphertext length...")
		}

		data, err := RSAOAEPDecrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("failed to rsa oaep decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("failed to produce valid decrypted bytes")
		}
	}
}
