package main

// Minimal implementation based on work by Satori:
// @link: https://github.com/satori/go.uuid

import (
	"encoding/hex"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// return a v4 uuid in standdard GUID string format
func GUID() string {
	var d [16]byte
	rand.Read(d[:])
	d[6] = (d[6] & 0x0f) | (4 << 4)
	d[8] = (d[8] & 0xbf) | 0x80
	b := make([]byte, 36)
	b[8], b[13], b[18], b[23] = '-', '-', '-', '-'
	hex.Encode(b[0:8], d[0:4])
	hex.Encode(b[9:13], d[4:6])
	hex.Encode(b[14:18], d[6:8])
	hex.Encode(b[19:23], d[8:10])
	hex.Encode(b[24:], d[10:])
	return string(b)
}
