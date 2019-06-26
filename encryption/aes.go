package main

// The CTR mode requires an IV, which becomes both part of the passed message,
// as well as used to "initialize" the CTR stream.  This means we cannot reuse
// the same CTR stream instance for decryption and encryption.
//
// Meanwhile, GCM works differently in that the AEAD structure can accept the
// nonce per operation, meaning it has less initializations since we can reuse
// the object.
//
// Therefore, I have a function to initialize the GCM AEAD structure
// independent of encryption and decryption.
//
// As an added bonus, it is possible to increment a nonce rather than
// randomly generate a new one per operation, which may impact the
// performance and thus a GCM function that accepts a nonce is available.

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Returns a AES GCM instance configured with the key.
func GCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return mode, nil
}

// Expects a nonce as a parameter and returns the combined ciphertext.
func GCMEncryptNonce(mode cipher.AEAD, nonce, message []byte) []byte {
	return append(nonce[:mode.NonceSize()], mode.Seal(nil, nonce[:mode.NonceSize()], message, nil)...)
}

// Generate crypto/rand nonce and return it with ciphertext.
func GCMEncrypt(mode cipher.AEAD, message []byte) ([]byte, error) {
	nonce := make([]byte, mode.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return GCMEncryptNonce(mode, nonce, message), nil
}

// decrypt using AEAD and expected nonce as part of ciphertext.
func GCMDecrypt(mode cipher.AEAD, ciphertext []byte) ([]byte, error) {
	return mode.Open(nil, ciphertext[:mode.NonceSize()], ciphertext[mode.NonceSize():], nil)
}

// Generate an iv and encrypt using CTR and return both as ciphertext.
func CTREncrypt(key, message []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, block.BlockSize()+len(message))

	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ciphertext, err
	}

	mode := cipher.NewCTR(block, iv)
	mode.XORKeyStream(ciphertext[block.BlockSize():], message)

	return ciphertext, nil
}

// Handles CTR decryption by extracting the iv from the ciphertext
func CTRDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	message := make([]byte, len(ciphertext)-block.BlockSize())
	mode := cipher.NewCTR(block, ciphertext[:block.BlockSize()])
	mode.XORKeyStream(message, ciphertext[block.BlockSize():])

	return message, nil
}
