package main

// a "simple" demonstration of parsing a token with
// permissions (eg. scope) embedded in the claims

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func main() {

	// set test case
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwicGVybXMiOlsiYWRtaW4iXX0.Vg0yGqo7kJysc_pTul7SWWCX5UGxnCWnaURXPPU1eaA"
	secret := []byte("secret")

	// parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrHashUnavailable
		}
		return secret, nil
	})

	// validate
	if err != nil || !token.Valid {
		fmt.Printf("invalid token: %s\n", err)
		return
	}
	fmt.Printf("token: %+v\n", token)

	// check claims for perms and extract
	var perms []interface{}
	if t, ok := token.Claims.(jwt.MapClaims); !ok {
		fmt.Println("failed to read claims...")
		return
	} else if p, ok := t["perms"]; !ok {
		fmt.Println("failed to get perms...")
		return
	} else {

		// @note: you cannot assume []string
		perms, ok = p.([]interface{})
		if !ok {
			fmt.Printf("unable to cast to read perms")
			return
		}
		fmt.Printf("perms: %+v\n", perms)
	}

	// cast each perm to check for admin
	for _, p := range perms {
		if s, ok := p.(string); ok {
			if s == "admin" {
				fmt.Println("Success!")
			}
		}
	}
}
