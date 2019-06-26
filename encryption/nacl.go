package main

// Lightweight nacl box abstraction.
//
// It is possible to use the same key pair, but most examples
// suggest exchanging public keys for normal communication.
//
// It seems nacl does not automate nonces.
//
// There is also the option for precomputation, which is worth
// demonstrating for the purposes of a benchmark, and is a feature
// not present or available to RSA.

import (
	"crypto/rand"
	"golang.org/x/crypto/nacl/box"
)

func NaClEncryptWithNonce(pub, priv *[32]byte, nonce *[24]byte, message []byte) []byte {
	return box.Seal(nonce[:], message, nonce, pub, priv)
}

func NaClEncrypt(pub, priv *[32]byte, message []byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nonce[:], err
	}
	return NaClEncryptWithNonce(pub, priv, &nonce, message), nil
}

func NaClDecrypt(pub, priv *[32]byte, ciphertext []byte) ([]byte, bool) {
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])
	return box.Open(nil, ciphertext[24:], &nonce, pub, priv)
}

func NaClPrecomputeEncryptWithNonce(key *[32]byte, nonce *[24]byte, message []byte) []byte {
	return box.SealAfterPrecomputation(nonce[:], message, nonce, key)
}

func NaClPrecomputeEncrypt(key *[32]byte, message []byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nonce[:], err
	}
	return NaClPrecomputeEncryptWithNonce(key, &nonce, message), nil
}

func NaClPrecomputeDecrypt(key *[32]byte, ciphertext []byte) ([]byte, bool) {
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])
	return box.OpenAfterPrecomputation(nil, ciphertext[24:], &nonce, key)
}
