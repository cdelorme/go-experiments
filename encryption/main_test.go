package main

// easily modified key and data sizes for manipulating performance comparisons
var KeySize int = 256 / 8                            // 32 bytes; AES key size
var RSAKeySize int = 256 * 8                         // 256 bytes
var MessageSize int = 508                            // 576 - 60 (IP Headers) - 8 (UDP Headers)
var GCMMessageSize int = MessageSize - 28            // 16 padding, 12 nonce
var RSAMessageSize int = 190                         // 190*8, 191 fails with 2048 key, too large?
var NaClMessageSize int = 190                        // to match RSA; no such size limits appear to apply here though
var CTRMessageSize int = MessageSize - 16            // 16 iv
var SignedCTRMessageSize int = MessageSize - 16 - 32 // 16 iv, 32 signature
var SignedMessageSize int = MessageSize - 32         // 32 signature
