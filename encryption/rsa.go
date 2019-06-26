package main

// RSA implementation abstraction, namely using OAEP padding.
//
// A short consideration was made regarding PSS, but it turns out
// that the use case involves both parties using key pairs and
// sharing their public keys.

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Leverage OAEP padding method to encrypt
func RSAOAEPEncrypt(key *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key, message, nil)
}

// Leverage OAEP padding method to decrypt
func RSAOAEPDecrypt(key *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, key, ciphertext, nil)
}
