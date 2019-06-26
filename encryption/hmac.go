package main

// Two simple functions to encapsulate the Sign and Verify processes.

import (
	"crypto/hmac"
	"crypto/sha256"
)

// The message must be no greater than 476 bytes.
func Sign(key, message []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// Use hmac.Equal to compare signatures.
func Verify(signature, key, message []byte) bool {
	return hmac.Equal(signature, Sign(key, message))
}
