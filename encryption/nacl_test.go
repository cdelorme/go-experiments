package main

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/nacl/box"
	"testing"
)

func TestNaCl(t *testing.T) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message[:]); err != nil {
		t.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate keypair: %s", err)
	}

	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		t.Fatalf("failed to generate nonce: %s", err)
	}

	ciphertext, err := NaClEncrypt(pub, priv, message)
	if err != nil {
		t.Fatalf("failed to encrypt with NaCl: %s", err)
	}

	// since the nonce is part of the data
	decrypted, ok := NaClDecrypt(pub, priv, ciphertext)
	if !ok {
		t.Fatal("failed to decrypt the message")
	} else if !bytes.Equal(decrypted, message) {
		t.Fatalf("Message does not match %s", decrypted)
	}
}

func TestNaClPrecompute(t *testing.T) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message[:]); err != nil {
		t.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate keypair: %s", err)
	}

	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		t.Fatalf("failed to generate nonce: %s", err)
	}

	var shared [32]byte
	box.Precompute(&shared, pub, priv)

	ciphertext, err := NaClPrecomputeEncrypt(&shared, message)
	if err != nil {
		t.Fatalf("failed to encrypt precomputed NaCl: %s", err)
	}

	// since the nonce is part of the data
	decrypted, ok := NaClPrecomputeDecrypt(&shared, ciphertext)
	if !ok {
		t.Fatal("failed to decrypt the message...")
	} else if !bytes.Equal(decrypted, message) {
		t.Fatalf("Message does not match %s", decrypted)
	}
}

///
// These test cases are more of a 1:1 comparison against RSA
///

func BenchmarkNaClEncrypt(b *testing.B) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("failed to generate keypair: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ciphertext, err := NaClEncrypt(pub, priv, message)
		if err != nil {
			b.Fatalf("failed to encrypt with NaCl: %s", err)
		} else if len(ciphertext) == 0 {
			b.Fatal("failed to produce a ciphertext...")
		}
	}
}

func BenchmarkNaClDecrypt(b *testing.B) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("failed to generate keypair: %s", err)
	}

	ciphertext, err := NaClEncrypt(pub, priv, message)
	if err != nil {
		b.Fatalf("failed to encrypt with NaCl: %s", err)
	} else if len(ciphertext) == 0 {
		b.Fatal("failed to produce a ciphertext...")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, ok := NaClDecrypt(pub, priv, ciphertext)
		if !ok {
			b.Fatal("failed to NaCl decrypt...")
		} else if !bytes.Equal(data, message) {
			b.Fatal("failed to produce valid decrypted bytes")
		}
	}
}

func BenchmarkNaCl(b *testing.B) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("failed to generate keypair: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ciphertext, err := NaClEncrypt(pub, priv, message)
		if err != nil {
			b.Fatalf("failed to encrypt with NaCl: %s", err)
		} else if len(ciphertext) == 0 {
			b.Fatal("failed to produce a ciphertext...")
		}

		data, ok := NaClDecrypt(pub, priv, ciphertext)
		if !ok {
			b.Fatal("failed to NaCl decrypt...")
		} else if !bytes.Equal(data, message) {
			b.Fatal("failed to produce valid decrypted bytes")
		}
	}
}

///
// precompute is expected to offer higher performance over RSA
///

func BenchmarkNaClPrecomputeEncrypt(b *testing.B) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("failed to generate keypair: %s", err)
	}

	var shared [32]byte
	box.Precompute(&shared, pub, priv)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ciphertext, err := NaClPrecomputeEncrypt(&shared, message)
		if err != nil {
			b.Fatalf("failed to encrypt with NaCl: %s", err)
		} else if len(ciphertext) == 0 {
			b.Fatal("failed to produce a ciphertext...")
		}
	}

}
func BenchmarkNaClPrecomputeDecrypt(b *testing.B) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("failed to generate keypair: %s", err)
	}

	var shared [32]byte
	box.Precompute(&shared, pub, priv)

	ciphertext, err := NaClPrecomputeEncrypt(&shared, message)
	if err != nil {
		b.Fatalf("failed to encrypt with NaCl: %s", err)
	} else if len(ciphertext) == 0 {
		b.Fatal("failed to produce a ciphertext...")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, ok := NaClPrecomputeDecrypt(&shared, ciphertext)
		if !ok {
			b.Fatal("failed to NaCl decrypt...")
		} else if !bytes.Equal(data, message) {
			b.Fatal("failed to produce valid decrypted bytes")
		}
	}
}
func BenchmarkNaClPrecompute(b *testing.B) {
	message := make([]byte, NaClMessageSize)
	if _, err := rand.Read(message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("failed to generate keypair: %s", err)
	}

	var shared [32]byte
	box.Precompute(&shared, pub, priv)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ciphertext, err := NaClPrecomputeEncrypt(&shared, message)
		if err != nil {
			b.Fatalf("failed to encrypt with NaCl: %s", err)
		} else if len(ciphertext) == 0 {
			b.Fatal("failed to produce a ciphertext...")
		}

		data, ok := NaClPrecomputeDecrypt(&shared, ciphertext)
		if !ok {
			b.Fatal("failed to NaCl decrypt...")
		} else if !bytes.Equal(data, message) {
			b.Fatal("failed to produce valid decrypted bytes")
		}
	}
}

///
// It is worth noting that we could write tests using incremental nonce
// which might put the precomputed performance very nearly on par with
// GCM encryption.
///
