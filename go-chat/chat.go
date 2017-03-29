package main

type Conn interface {
	Send(Message)
}

type Chat struct {
	connections []Conn
}

func (self *Chat) Add(c Conn) {
	self.connections = append(self.connections, c)
}

func (self *Chat) Send(m Message) {
	for i := range self.connections {
		self.connections[i].Send(m)
	}
}

type ChannelChat struct {
	channel chan Message
	Chat
}

func (self *ChannelChat) process() {
	for m := range self.channel {
		self.Chat.Send(m)
	}
}

func (self *ChannelChat) Add(c Conn) {
	if self.channel == nil {
		self.channel = make(chan Message)
		go self.process()
	}
	self.Chat.Add(c)
}

func (self *ChannelChat) Send(m Message) {
	self.channel <- m
}
