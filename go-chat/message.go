package main

import (
	"strconv"
	"time"
)

type Message struct {
	Content   string
	Author    string
	Timestamp time.Time
}

func (self *Message) Byte() []byte {
	return []byte(strconv.Itoa(len(self.Content)) + "\n" + self.Content + "\n" + strconv.Itoa(len(self.Author)) + "\n" + self.Author + "\n" + strconv.FormatInt(self.Timestamp.UnixNano(), 10))
}
