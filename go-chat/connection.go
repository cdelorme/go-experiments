package main

import (
	"io"
	"sync"
)

type Connection struct {
	stream io.Writer
}

func (self *Connection) Send(m Message) {
	self.stream.Write(m.Byte())
}

type ConnectionMutex struct {
	sync.Mutex
	Connection
}

func (self *ConnectionMutex) Send(m Message) {
	self.Lock()
	self.Connection.Send(m)
	self.Unlock()
}
