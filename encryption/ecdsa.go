package main

// On its own ECDSA does not provide encryption and decryption, but rather
// signing and verification.
//
// Thus we have to add layers to implement encryption and decryption.

// import (
// 	"crypto"
// 	"crypto/ecdsa"
// )

// func ECDSASign(key *crypto.PublicKey, message []byte) ([]byte, error) {

// 	return []byte{}, nil
// }

// func ECDSAVerify(key *ecdsa.PrivateKey, cyphertext []byte) ([]byte, error) {

// 	return []byte{}, nil
// }

// ECDSA implementation abstraction

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/rand"
// 	// "crypto/tls"
// 	// "crypto/x509"
// 	// "crypto/x509/pkix"
// 	// "fmt"
// 	// "math/big"
// 	// "os"
// 	"testing"
// 	// "time"
// )

// // Create a keypair dynamically.
// //
// // Use the public key to encrypt.
// //
// // Use the private key to decrypt.
// //
// // Use the private key to encrypt.
// //
// // Use the public key to decrypt.
// func TestRSA(t *testing.T) {
// 	priv, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
// 	if err != nil {
// 		t.Fatalf("failed to create a private key: %s", err)
// 	}
// 	t.Logf("Public Key: %x", priv.PublicKey)

// }

///
///
///

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/md5"
// 	"crypto/rand"
// 	"fmt"
// 	"hash"
// 	"io"
// 	"math/big"
// 	"os"
// )

// func main() {

// 	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

// 	privatekey := new(ecdsa.PrivateKey)
// 	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	var pubkey ecdsa.PublicKey
// 	pubkey = privatekey.PublicKey

// 	fmt.Println("Private Key :")
// 	fmt.Printf("%x \n", privatekey)

// 	fmt.Println("Public Key :")
// 	fmt.Printf("%x \n", pubkey)

// 	// Sign ecdsa style

// 	var h hash.Hash
// 	h = md5.New()
// 	r := big.NewInt(0)
// 	s := big.NewInt(0)

// 	io.WriteString(h, "This is a message to be signed and verified by ECDSA!")
// 	signhash := h.Sum(nil)

// 	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signhash)
// 	if serr != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	signature := r.Bytes()
// 	signature = append(signature, s.Bytes()...)

// 	fmt.Printf("Signature : %x\n", signature)

// 	// Verify
// 	verifystatus := ecdsa.Verify(&pubkey, signhash, r, s)
// 	fmt.Println(verifystatus) // should be true
// }
