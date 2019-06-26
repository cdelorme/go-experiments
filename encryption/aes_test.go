package main

// Test cases and performance comparison between AES GCM, AES CTR.
//
// Two AES GCM solutions, one uses kernel random per nonce with an
// assumed safety limit of 2^32 before potential redundancy, and the
// other uses uint64 casting with an exponentially higher maximum.
//
// Two forms of AES CTR, one with signatures one without; however
// there are absolutely no scenarios where not signing is acceptable.

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

// A test of encryption and decryption using GCM.
func TestGCM(t *testing.T) {
	key := make([]byte, KeySize)
	message := make([]byte, GCMMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		t.Fatalf("failed to create message: %s", err)
	}

	mode, err := GCM(key)
	if err != nil {
		t.Fatalf("failed to create GCM AEAD: %s", err)
	}

	ciphertext, err := GCMEncrypt(mode, message)
	if err != nil {
		t.Fatalf("failed to gcm encrypt: %s", err)
	}

	t.Logf("Message Size: %d", len(message))
	t.Logf("Ciphertext Size: %d", len(ciphertext))

	data, err := GCMDecrypt(mode, ciphertext)
	if err != nil {
		t.Fatalf("failed to gcm decrypt: %s", err)
	}

	if !bytes.Equal(data, message) {
		t.Fatal("decrypted message does not equal original...")
	}
}

// A test of encryption and decryption using CTR.
func TestCTR(t *testing.T) {
	key := make([]byte, KeySize)
	message := make([]byte, CTRMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		t.Fatalf("failed to create message: %s", err)
	}

	ciphertext, err := CTREncrypt(key, message)
	if err != nil {
		t.Fatalf("failed to ctr encrypt: %s", err)
	}

	t.Logf("Message Size: %d", len(message))
	t.Logf("Cipher Size: %d", len(ciphertext))

	data, err := CTRDecrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("failed to ctr decrypt: %s", err)
	}

	if !bytes.Equal(data, message) {
		t.Fatal("decrypted message does not equal original...")
	}
}

func BenchmarkAESGCMEncrypt(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, GCMMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}
	mode, err := GCM(key)
	if err != nil {
		b.Fatalf("failed to prepare gcm: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := GCMEncrypt(mode, message)
		if err != nil {
			b.Fatalf("failed to gcm encrypt: %s", err)
		} else if len(ciphertext) != len(message)+mode.NonceSize()+mode.Overhead() {
			b.Fatal("invalid ciphertext size...")
		}
	}
}

func BenchmarkAESGCMEncryptIncNonce(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, GCMMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}
	mode, err := GCM(key)
	if err != nil {
		b.Fatalf("failed to prepare gcm: %s", err)
	}
	nonce := make([]byte, mode.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		b.Fatalf("failed to generate nonce: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var inonce uint64 = uint64(nonce[0]) | uint64(nonce[1])<<8 | uint64(nonce[2])>>16 | uint64(nonce[3])<<24 | uint64(nonce[4])<<32 | uint64(nonce[5])<<40 | uint64(nonce[6])<<48 | uint64(nonce[7])<<56
		inonce += 1
		nonce[0], nonce[1], nonce[2], nonce[3], nonce[4], nonce[5], nonce[6], nonce[7] = byte(inonce), byte(inonce>>8), byte(inonce>>16), byte(inonce>>24), byte(inonce>>32), byte(inonce>>40), byte(inonce>>48), byte(inonce>>56)
		ciphertext := GCMEncryptNonce(mode, nonce, message)
		if len(ciphertext) != len(message)+mode.NonceSize()+mode.Overhead() {
			b.Fatal("invalid ciphertext size...")
		}
	}
}

func BenchmarkAESGCMDecrypt(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, GCMMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}
	mode, err := GCM(key)
	if err != nil {
		b.Fatalf("failed to prepare gcm: %s", err)
	}
	ciphertext, err := GCMEncrypt(mode, message)
	if err != nil {
		b.Fatalf("failed to gcm encrypt: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		if data, err := GCMDecrypt(mode, ciphertext); err != nil {
			b.Fatalf("failed to gcm decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("decryption yielded invalid bytes...")
		}
	}
}

func BenchmarkAESGCMFull(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, GCMMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}
	mode, err := GCM(key)
	if err != nil {
		b.Fatalf("failed to prepare gcm: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := GCMEncrypt(mode, message)
		if err != nil {
			b.Fatalf("failed to gcm encrypt: %s", err)
		} else if len(ciphertext) != len(message)+mode.NonceSize()+mode.Overhead() {
			b.Fatal("invalid ciphertext size...")
		} else if data, err := GCMDecrypt(mode, ciphertext); err != nil {
			b.Fatalf("failed to gcm decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("decryption yielded invalid bytes...")
		}
	}
}

func BenchmarkAESGCMFullIncNonce(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, GCMMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}
	mode, err := GCM(key)
	if err != nil {
		b.Fatalf("failed to prepare gcm: %s", err)
	}
	nonce := make([]byte, mode.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		b.Fatalf("failed to generate nonce: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var inonce uint64 = uint64(nonce[0]) | uint64(nonce[1])<<8 | uint64(nonce[2])>>16 | uint64(nonce[3])<<24 | uint64(nonce[4])<<32 | uint64(nonce[5])<<40 | uint64(nonce[6])<<48 | uint64(nonce[7])<<56
		inonce += 1
		nonce[0], nonce[1], nonce[2], nonce[3], nonce[4], nonce[5], nonce[6], nonce[7] = byte(inonce), byte(inonce>>8), byte(inonce>>16), byte(inonce>>24), byte(inonce>>32), byte(inonce>>40), byte(inonce>>48), byte(inonce>>56)
		ciphertext := GCMEncryptNonce(mode, nonce, message)
		if len(ciphertext) != len(message)+mode.NonceSize()+mode.Overhead() {
			b.Fatal("invalid ciphertext size...")
		} else if data, err := GCMDecrypt(mode, ciphertext); err != nil {
			b.Fatalf("failed to gcm decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("decryption yielded invalid bytes...")
		}
	}
}

func BenchmarkAESCTREncrypt(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, CTRMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := CTREncrypt(key, message)
		if err != nil {
			b.Fatalf("failed to ctr encrypt: %s", err)
		} else if len(ciphertext) <= len(message) {
			b.Fatal("wat?")
		}
	}
}

func BenchmarkAESCTRDecrypt(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, CTRMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}
	ciphertext, err := CTREncrypt(key, message)
	if err != nil {
		b.Fatalf("failed to ctr encrypt: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := CTRDecrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("failed to ctr decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("decrypted bytes do not match original...")
		}
	}
}

func BenchmarkAESCTRFull(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, CTRMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := CTREncrypt(key, message)
		if err != nil {
			b.Fatalf("failed to ctr encrypt: %s", err)
		}
		data, err := CTRDecrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("failed to ctr decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("decrypted bytes do not match original...")
		}
	}
}

func BenchmarkAESCTRHMACEncrypt(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, SignedCTRMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := CTREncrypt(key, message)
		if err != nil {
			b.Fatalf("failed to ctr encrypt: %s", err)
		} else if len(ciphertext) <= len(message) {
			b.Fatal("wat?")
		}
		signature := Sign(key, ciphertext)
		if len(signature) != 32 {
			b.Fatal("failed to generate HMAC signature...")
		}
	}
}

func BenchmarkAESCTRHMACDecrypt(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, SignedCTRMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}
	ciphertext, err := CTREncrypt(key, message)
	if err != nil {
		b.Fatalf("failed to ctr encrypt: %s", err)
	}
	signature := Sign(key, ciphertext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !Verify(signature, key, ciphertext) {
			b.Fatal("signatured failed verification...")
		}
		data, err := CTRDecrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("failed to ctr decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("decrypted bytes do not match original...")
		}
	}
}

func BenchmarkAESCTRHMACFull(b *testing.B) {
	key := make([]byte, KeySize)
	message := make([]byte, SignedCTRMessageSize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		b.Fatalf("failed to create key: %s", err)
	} else if _, err := io.ReadFull(rand.Reader, message); err != nil {
		b.Fatalf("failed to create message: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := CTREncrypt(key, message)
		if err != nil {
			b.Fatalf("failed to ctr encrypt: %s", err)
		}
		signature := Sign(key, ciphertext)
		if !Verify(signature, key, ciphertext) {
			b.Fatal("signatured failed verification...")
		}
		data, err := CTRDecrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("failed to ctr decrypt: %s", err)
		} else if !bytes.Equal(data, message) {
			b.Fatal("decrypted bytes do not match original...")
		}
	}
}
