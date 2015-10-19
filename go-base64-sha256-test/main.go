package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// a simple test to verify the differences between
// sha256 and base64 for possible bugs in token creation

type Thing struct {
	Message     string    `json:"message,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}

func main() {

	secret := "a secret shared encryption code"
	message := Thing{
		Message: "Dangajits",
		Date:    time.Now(),
	}

	jsonMessage, _ := json.Marshal(message)

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(jsonMessage))
	shaout := h.Sum(nil)

	b64m := base64.StdEncoding.EncodeToString(jsonMessage)
	b64um := base64.URLEncoding.EncodeToString(jsonMessage)
	b64jm := base64.StdEncoding.EncodeToString(shaout)
	b64jum := base64.URLEncoding.EncodeToString(shaout)

	fmt.Printf("%+v\n", message)
	fmt.Printf("%s\n", jsonMessage)
	fmt.Println(b64m)
	fmt.Println(b64um)
	fmt.Println(b64jm)
	fmt.Println(b64jum)

	// the lesson?
	// sha256 always produces characters that are url friendly
	// jwt token signature encoding can use stdencoding
	// however, the payload is just basic json
	// so that needs urlencoding

	// the equal signs at the end are the result of padding
	// which can be removed using WithPadding(NoPadding)
	// @link: https://golang.org/pkg/encoding/base64/
}
