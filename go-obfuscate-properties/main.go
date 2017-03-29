package main

// how to obfuscate properties
// such that prints will not
// display them, using my
// logger as a printer to
// show that it works beyond
// simple fmt commands

import (
	"github.com/cdelorme/go-log"
	"strconv"
	// "fmt"
)

// type Password string

// func (self *Password) GoString() string {
// 	return ""
// }

type Configuration struct {
	Port     int
	Address  string
	Username string
	// Password Password
	Password string
}

// @note: sadly there is no simpler solution
//        short of triple json conversion to
//        map where you can drop properties
func (self *Configuration) MarshalJSON() ([]byte, error) {
	return []byte(`{"port": ` + strconv.Itoa(self.Port) + `, "address": "` + self.Address + `", "username": "` + self.Username + `"}`), nil
}

// @note: this is used by `%s` and `%v`
func (self Configuration) String() string {
	s, _ := self.MarshalJSON()
	return string(s)
}

// @note: this is used by `%#v`
func (self Configuration) GoString() string {
	return self.String()
}

func main() {
	l := log.Logger{}
	c := Configuration{Port: 3000, Address: "Nope", Username: "Narp", Password: "u_pick_it"}

	// @note: the password property should be absent from all of these
	l.Error("String: %s", c)
	l.Error("Verbose: %v", c)
	l.Error("Very Verbose: %+v", c)
	l.Error("GoString: %#v", c)
}
