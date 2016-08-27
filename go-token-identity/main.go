package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const tokenString = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJFbGx1Y2lhbiBDbG91ZCBQbGF0Zm9ybSIsInN1YiI6ImJlbi52aXRhbGVAZWxsdWNpYW4ubWUiLCJ0ZW5hbnQiOnsiaWQiOiI1M2IyZjFkMDY0MWEwYTBmMDAwYzJmMTEiLCJhbGlhcyI6ImVsbHVjaWFuIn0sImlhdCI6MTQzNzQxNjY1MH0.iYhjy8f7s3_dsdSttyc1kX9E2vAe9b2LFwkDIlwI1vo"

var JwtSharedPassword = []byte("it scales right up")
var errorBadSigningMethod = errors.New("Unexpected signing method")

type cache struct {
	sync.Mutex
	counts map[string]int
	second int
}

func (self *cache) Next(tenant, uri string) int {
	self.Lock()
	defer self.Unlock()
	if self.counts == nil {
		self.counts = map[string]int{}
	}
	if time.Now().Second() > self.second {
		delete(self.counts, tenant+uri+strconv.Itoa(self.second))
		self.second = time.Now().Second()
	}
	self.counts[tenant+uri+strconv.Itoa(self.second)] += 1
	return self.counts[tenant+uri+strconv.Itoa(self.second)]
}

func tenant(token *jwt.Token) string {

	// attempt to extract the tenant
	t, ok := token.Claims["tenant"]
	if !ok {
		fmt.Printf("Unable to grab tenant id from token (%v)", token)
		return ""
	}

	// attempt to cast the tenant to a map
	d, ok := t.(map[string]interface{})
	if !ok {
		fmt.Printf("Unable to grab tenant id from token (%v)", token)
		return ""
	}

	// attempt to extract the id
	id, ok := d["id"]
	if !ok {
		fmt.Printf("Unable to grab tenant id from token (%v)", token)
		return ""
	}

	// attempt to convert the id to a string
	sid, ok := id.(string)
	if !ok {
		fmt.Printf("Unable to grab tenant id from token (%v)", token)
		return ""
	}

	return sid
}

func main() {

	// // create cache instance
	// c := cache{}

	// // test next()
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))

	// // sleep for 2 seconds
	// time.Sleep(time.Second * 2)

	// // try again
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))
	// fmt.Printf("%d\n", c.Next("fuck", "stick"))

	// create a token from a string
	token, e := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorBadSigningMethod
		}

		return JwtSharedPassword, nil
	})
	fmt.Printf("Token: %#v\n", token)
	fmt.Printf("Error: %#v\n", e)

	// try extracting the tenantid from it
	fmt.Printf("Tenant: %s\n", tenant(token))
}
