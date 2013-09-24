package bitmask

import (
	"fmt"
	"net"
)

type Peer struct {
	socket  net.Conn
	context Context
}

func NewPeer(socket net.Conn, context *Context) *Peer {
	this := new(Peer)
	this.socket = socket
	this.context = context
	go this.run()

	return this
}

func (this *Peer) run() {
	defer this.socket.Close()
	fmt.Println("New user connected")
}
