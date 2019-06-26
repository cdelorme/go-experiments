package main

// Benchmarks to compare the performance of Sign and Verify.
//
// There is no "full" benchmark, because Verify is already performing Sign.

import (
	"crypto/rand"
	"io"
	"testing"
)

// This verifies both Sign & Verify.
func TestHMAC(t *testing.T) {
	key := make([]byte, KeySize)
	message := make([]byte, SignedMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Fatalf("failed to generate key: %s", err)
	} else if _, err = io.ReadFull(rand.Reader, message); err != nil {
		t.Fatalf("failed to generate message: %s", err)
	} else if !Verify(Sign(key, message), key, message) {
		t.Fatal("failed to verify own signature...")
	}
}

func BenchmarkHMACSign(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, SignedMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to generate key: %s", err)
	} else if _, err = io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to generate message: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signature := Sign(key, message)
		if len(signature) != 32 {
			b.Fatalf("expected 32 bytes, got %d...", len(signature))
		}
	}
}

func BenchmarkHMACVerify(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, SignedMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to generate key: %s", err)
	} else if _, err = io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to generate message: %s", err)
	}
	signature := Sign(key, message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !Verify(signature, key, message) {
			b.Fatal("failed to verify own sginature...")
		}
	}
}
