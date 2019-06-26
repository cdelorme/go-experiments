
# encryption

This is an overarching demonstration of packages that enable encryption, when you may want to use them, and how they perform.

I will finish with an overview of concerns related to encryption technology.


## types

There are two primary types of encryption:

- asymmetric
- symmetric

They cannot be compared.

They serve different purposes.

They are often used together.


### asymmetric

The term asymmetric is used to describe of a public and private key to separate encryption and decryption in interesting ways.

**The private key should never be shared.**  The holder of the private key can decrypt any data that has been encrypted with the public key.  It can also be used to encrypt or more commonly **sign** data that the public key can they verify.

The public key can be also be used to decrypt data encrypted with the private key, but is more often used to verify a signature.  _A valid signature means that the data most certainly came from the holder of the private key._

The "gold standard" for asymmetric encryption is RSA.

A modern and comparable technology is ECC (elliptic curve encryptiun), which relies on theoretical math and offers smaller key pairs with equal levels of security.

While experts like Colin Percival suggest to "**probably** avoid elliptical curves", it is perhaps subject to interpretation whether he meant implementing the elliptical curve algorithms themselves or the technology as a whole.

On the other hand relative large businesses like [CloudFlare](https://blog.cloudflare.com/ecdsa-the-digital-signature-algorithm-of-a-better-internet/) have embraced ECC due to the significant performance gains both at the server and each individual client, _particularly of interest for embedded and mobile devices_.


### symmetric

Symmetric encryption uses a single shared secret for both encryption and decryption.

It is important to understand that with symmetric encryption there are three aspects.

1. Algorithm
2. Block Cipher Mode
3. Authentication

Among algorithms the "gold standard" is AES, although there are others such as (3)DES and Blowfish many are compromised or proprietary.

There are many block cipher modes, common ones include CBC, CTR, and GCM.

While CBC is relatively dated, GCM and CTR are popular.

However, GCM features "authenticated encryption with associated data" (AEAD), which is the equivalent of the third aspect (authentication) embedded within the second.

**It is worth noting that Colin Percival has made two claims against GCM.**

1. GCM may be subject to side-channel attacks, _but this may be a matter of implementation_.
2. Separate signature validation may alleviate DOS attacks by discarding prior to decryption.

The final aspect, authentication, as a stand-alone solution is commonly implemented using "hash-based message authentication code" (HMAC) and may be referred to as a "signature".  _This is considered a necessary extra step for all block cipher modes that are not AEAD._


## performance

With both RSA and ECC, the following data can be compared:

- Key Generation
- Public Key Encryption
- Public Key Decryption
- Private Key Decryption
- Private Key Encryption

While RSA has giant keys it is considered a proven technology, whereas elliptic curves are still considered a young technology and there are concerns over unknown weaknesses in new algorithms and the patents on or creators of existing binary curves.


Both AES-GCM and AES-CTR+HMAC-SHA256 may be compared:

- Encryption
- Decryption

_It may be worth comparing both complexity of implementation as the block cipher modes used may not always incorporate signature generation, making it extra work that has to be done and which may also offset the final performance._

There are claims that block ciphers that handle signature generation are insecure, but similar claims can be made towards the ease of human error when handling signatures independently.


### benchmarks

Feel free to try this on your own system:

	go test -v -run=X -benchtime=10s -benchmem -bench=.

_You may have to set5 `-timeout=0` or a value beyond the 10m default, as these benchmarks take a while._

My results were as follows:

	goos: linux
	goarch: amd64
	pkg: github.com/cdelorme/go-experiments/encryption
	BenchmarkAESGCMEncrypt-12            	20000000	      1069 ns/op	    1040 B/op	       3 allocs/op
	BenchmarkAESGCMEncryptIncNonce-12    	30000000	       413 ns/op	    1024 B/op	       2 allocs/op
	BenchmarkAESGCMDecrypt-12            	50000000	       283 ns/op	     480 B/op	       1 allocs/op
	BenchmarkAESGCMFull-12               	10000000	      1365 ns/op	    1520 B/op	       4 allocs/op
	BenchmarkAESGCMFullIncNonce-12       	20000000	       832 ns/op	    1504 B/op	       3 allocs/op
	BenchmarkAESCTREncrypt-12            	10000000	      1898 ns/op	    1696 B/op	       8 allocs/op
	BenchmarkAESCTRDecrypt-12            	10000000	      1513 ns/op	    1696 B/op	       8 allocs/op
	BenchmarkAESCTRFull-12               	 5000000	      3055 ns/op	    3392 B/op	      16 allocs/op
	BenchmarkAESCTRHMACEncrypt-12        	 5000000	      3698 ns/op	    2176 B/op	      14 allocs/op
	BenchmarkAESCTRHMACDecrypt-12        	 5000000	      3109 ns/op	    2176 B/op	      14 allocs/op
	BenchmarkAESCTRHMACFull-12           	 2000000	      6739 ns/op	    4352 B/op	      28 allocs/op
	BenchmarkHMACSign-12                 	10000000	      1705 ns/op	     512 B/op	       6 allocs/op
	BenchmarkHMACVerify-12               	10000000	      1817 ns/op	     512 B/op	       6 allocs/op
	BenchmarkNaClEncrypt-12              	  500000	     33645 ns/op	     272 B/op	       2 allocs/op
	BenchmarkNaClDecrypt-12              	  500000	     33209 ns/op	     192 B/op	       1 allocs/op
	BenchmarkNaCl-12                     	  200000	     65210 ns/op	     464 B/op	       3 allocs/op
	BenchmarkNaClPrecomputeEncrypt-12    	10000000	      1258 ns/op	     272 B/op	       2 allocs/op
	BenchmarkNaClPrecomputeDecrypt-12    	20000000	       617 ns/op	     192 B/op	       1 allocs/op
	BenchmarkNaClPrecompute-12           	10000000	      1849 ns/op	     464 B/op	       3 allocs/op
	BenchmarkRSAEncrypt-12               	  300000	     50837 ns/op	   13625 B/op	      33 allocs/op
	BenchmarkRSADecrypt-12               	   10000	   1116908 ns/op	   34689 B/op	     123 allocs/op
	BenchmarkRSA-12                      	   10000	   1204564 ns/op	   48317 B/op	     156 allocs/op
	PASS
	ok  	github.com/cdelorme/go-experiments/encryption	384.009s


## concerns

I want to cover several loosely related topics in greater details:

- Reserved Bytes
- ECC (Elliptic Curve Encryption)
- Forward Secrecy (eg. Diffie-Hellman)
- PKI (Public Key Infrastructure)
- MITM (Man In The Middle)
- Side Channel (Timing)


## Reserved Bytes

The most common advice is not to write your own encryption, and so the majority of networked services will be written simply using HTTPS/TLS, and none of this will matter.

The exception is UDP, where you cannot leverage HTTPS/TLS encryption.

Therefore the majority of tests are written with UDP packet sizes in mind specifically to avoid fragmentation which exponentially increases failure due to lost packets.

This made a strong case for understanding the exact costs of each applied technology, where things like block size, padding, nonce, and signatures become relative and important:

- AES has a block size of 16.
- AES requires 16 bytes for padding.
- GCM requires 12 bytes for a nonce.
- CTR requires 16 bytes for the IV.
- An HMAC signature is 32 bytes.
- RSA 2048 has a block size of 256
- RSA 4096 has a block size of 512
- OAEP padding consumes 66 bytes.
- NaCl requires a 24 byte nonce.
- NaCl adds a 16 byte signature.

Therefore we can conclude that:

- AES-GCM costs 28 bytes
- AES-CTR+HMAC-SHA256 costs 48 bytes
- RSA2048+OEAP has a message size of 190 bytes
- RSA4096+OEAP has a message size of 446 bytes

While AES-GCM has an extra 20 bytes, the difference is obscured by the 16 byte block size.  Subtracting the 28 bytes of padding leaves you with exactly 30 blocks or 480 bytes of space, but the 48 bytes consumed by the AES-CTR+HMAC-SHA256 does not align well and leaves you with an 14 spare bytes and only 28 blocks or 448 bytes of message space.

With RSA the case is unique, as it seems the OAEP padding is treated as part of the block, hence the max message size is the block size minus the padding.

Unfortunately there is no way to ensure that RSA 4096 does not fragment, although it would only be two packets at most, and only occur during initial connections, which probably make the trade-off worth the added security.

In either case, RSA should only be used to exchange a symmetric key or discrete logarithm equivalent such as Diffie-Hellman.

As far as I can tell there is no "block size" for NaCl, it uses the message size verbatim, adding the 24 byte nonce, and 16 bytes which I believe are the signature.  _I only tested some small message sizes, but even with odd sized messages it produced the same results._


### ECC

Since ECC is based on complex theoretical mathematics there is still wide concern that there may be unknown vulnerabilities, which is why the use of standardized curves is recommended.

However many are rightly concerned that the common NIST (P224, P256, P384, P521) algorithms were created by the NSA, a branch of the US government.  _There is some history where they **may** have been aware of vulnerabilities in older implementations._

Most recommend using Curve25519/ED25519, but there does not appear to be a standard library option for this with golang.


### Forward Secrecy

The concept of forward secrecy applies when dealing with key exchange compromises.

In theory a discrete logarithm could be used in place of key pairs when establishing a connection, as the data exchanged does not compromise the secret.  _However, a certificate authority may still be required to avoid MITM proxy attacks._

It would probably be best to simply look up the Diffie-Hellman discrete logarithm, but suffice to say the basis of ECC works around the same concept.

Two parties pick a random value, agree upon portions of the formula, exchange the augmented version of their half of the puzzle, and are able to use the others augmented half via an elliptical curve to compute the same secret.  Any MITM listener only has the augmented halves, which means they cannot get the same answer as either of the two parties.


### PKI

A "public key infrastructure" (PKI) is summarily the method used to distribute public keys, which makes it the cornerstone that enables trust with encryption.

_It can literally be as simple as "included with the client"._

The majority of the internet relies on X509 certificate formats to establish a "circle of trust".

This works by having a third party called a "certificate authority" use their private key to add a signature to your X509 certificate.

From here a web browser or similar client would be responsible for checking that signature with the public key of that certificate authority.

_Unfortunately this means a compromised, corrupt, or malicious certificate authority could provide invalid certificates._

While there have been no CA private keys reported compromised, there have been social engineering and other types of compromises that did not rely on technology leading to certificates being issued to entities that did not own domains.

There does not appear to be regulations surrounding CAs, at the very least any regulations which can be enforced.  Instead, browser manufacturers like Google have gotten involved by removing certificate authorities who had shown a lack of compliance with "expectations".


### MITM

A "Man in the Middle" (MITM) attack is what certificate systems exist to guard against.

A traditional MITM attack on **unencrypted** data merely needs to intercept traffic to see what is happening.

When you add encryption, the MITM attack needs not only to intercept the messages but also the power to manipulate them or the knowledge _and processing power_ to decrypt them.

**It is worth pointing out that MITM attacks are not time sensitive; _they can collect traffic and store it, and in the future they may be able to decrypt it._**

These types of attacks are referred to as "passive" in that they are impossible to detect because they neither redirect nor alter the traffic.

The focal point of MITM attacks against encrypted traffic is the initial handshake between two systems.  _This is when the cipher and key are chosen for encryption are chosen._

To protect the exchange of cipher and key, a handshake process that involves an asymmetric key is used.  _A server shares its public key, and the cipher and key are encrypted such that only that servers private key may decrypt the message._

**This is where a MITM attack levels up so to speak; because it can intercept and proxy traffic and perform its own handshake with the actual server then supply its own key pair handshake to all clients.**

This form of interception is generally considered an "active" attack and is able to be detected as well as prevented.  However, the nature of the internet offers no such guarantees that everyone will choose to use secure configuration.

The only reliable way to protect against this form of MITM attack is to use a trusted third party to validate the servers key/certificate.  _This is where the PKI using certificate authorities as a circle-of-trust comes into play, and where key pairs transform into "certificate chains"._

A certificate authority will issue a verified certificate, and the client can check that certificate preventing the MITM from faking the handshake and generally _manipulating or decrypting_ communication.

**Even then this is still subject to many forms of attack!**

1. A MITM attack may use a compromised private key to hijack and proxy.
2. A MITM attack may be issued a private key from a compromised CA to hijack and proxy.
3. All interactions can be recorded in hopes that a key may be compromised in the future; _this leads to the need for forward secrecy._
4. A compromised client machine could be to told to trust a non-standard CA used to issue fake certificates.

We can eliminate the third concern by using a cipher with forward secrecy such as a discrete logarithm.

Choosing not to use an algorithm with forward secrecy can lead to all messages being recorded and decrypted in the event the private key is compromised in the future, _even if that key is no longer actively being used._

The fourth is a possibility regardless of the scheme used, though it may be easier at that point for the attacker to collect traffic from the victims machine directly.

A final noteworthy consideration is MITM efficacy against real-time applications.  In the event that the data being transmitted is time sensitive, even a real-time MITM proxy decryption is going to add latency and it will also be very expensive for the attacker.  Further, that data may loose relevance or value over time eliminating the value of forward secrecy.


### Side Channel

A side channel attack is one that uses operational times to glean information from the encryption implementation.

A naive implementation of a padding algorithm may be subject to such an attack.

By manipulating individual bytes and comparing time-to-failure it is possible to break an encryption scheme.

This is of course dependent on the implementation not using equal time steps for comparisons.


# conclusions

**Just use NaCl.**

No really.

It provides asymmetric (box) and symmetric (secretbox) support.

It transparently converts asymmetric to symmetric using a Diffie-Hellman like exchange of A/B (public/private keys).

It suffers none of the RSA encryption byte limitations.

If used raw its asymmetric performance is double that of RSA.

If used precomputed it competes with AES-GCM, beating AES-CTR+HMAC-SHA256.

It leverages ECC Curve25519 which is the most commonly supported and least drama-heavy curve (_eg. not RSA NIST_).

It's very easy to implement without introducing accidental risks; except nonce all risky parts are abstracted.

If you only have AES support then AES-GCM is preferred and with incremental nonce is an order of magnitude faster than AES-CTR+HMAC-SHA256.  _DOS attack vectors are still relevant so there may be a valid use case for separated HMAC over AEAD._

If you must use RSA be aware of the abysmal decryption performance as it will greatly reduce your servers throughput.


# references

- [RSA Examples](https://l1z2g9.github.io/2016/11/04/RSA-Encrypt-Decrypt-with-Golang/)
- [rsa example](https://github.com/brainattica/Golang-RSA-sample/blob/master/rsa_sample.go)
- [ECDSA: The digital signature algorithm of a better internet](https://blog.cloudflare.com/ecdsa-the-digital-signature-algorithm-of-a-better-internet/)
- [Cryptographic Right Answers](http://www.daemonology.net/blog/2009-06-11-cryptographic-right-answers.html)
- [Encrypt-then-MAC](http://www.daemonology.net/blog/2009-06-24-encrypt-then-mac.html)
- [scrypt](https://godoc.org/golang.org/x/crypto/scrypt)
- [practical cryptography with go (using nacl /w secretbox)](https://leanpub.com/gocrypto/read)
- [safe curves](http://safecurves.cr.yp.to/)
- [AES-256 GCM Encryption Example in Golang](https://gist.github.com/kkirsche/e28da6754c39d5e7ea10)
- [Authenticated Encryption with Additional Data using AES-GCM](https://download.libsodium.org/doc/secret-key_cryptography/aes-256-gcm.html)
- [Encrypting and decrypting data](https://astaxie.gitbooks.io/build-web-application-with-golang/en/09.6.html)
- [Everything you need to know about cryptography in 1 hour - Colin Percival](https://www.youtube.com/watch?v=jzY3m5Kv7Y8)
- [Everything you need to know about cryptography in 1 hour](http://www.daemonology.net/papers/crypto1hr.pdf)
- [ECC Curve 25519](https://godoc.org/golang.org/x/crypto/ed25519)
- [Curve 25519](https://godoc.org/golang.org/x/crypto/curve25519)
- [fun project](https://www.reddit.com/r/golang/comments/7n0z01/aesctr_with_hmac_sha256/)
- [ecdsa example](https://www.socketloop.com/tutorials/golang-example-for-ecdsa-elliptic-curve-digital-signature-algorithm-functions)
- [Considerations when using AES-GCM for encrypting files](https://blog.secure-monkey.com/?p=94)
- [golang x/crypto/curve25519](https://godoc.org/golang.org/x/crypto/curve25519)
- [golang x/crypto/ed25519](https://godoc.org/golang.org/x/crypto/ed25519)
- [Go Lang NACL Cryptography](https://8gwifi.org/docs/go-nacl.jsp)
- [golang server](https://austburn.me/blog/golang-server.html)
- [Cryptographic Best Practices](https://gist.github.com/atoponce/07d8d4c833873be2f68c34f9afc5a78a)
