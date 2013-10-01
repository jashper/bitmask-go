package bitmask

import (
	"fmt"
	"net"
)

type Network struct {
	context   *Context
	listeners []*Listener
}

type Listener struct {
	context *Context
	socket  net.Listener
}

func NewNetwork(context *Context) (this *Network) {
	this = new(Network)
	this.context = context
	this.listeners = make([]*Listener, 0)
	this.SpawnListener(":7001")

	return this
}

func NewListener(context *Context, port string) (this *Listener, success bool) {
	this = new(Listener)
	this.context = context

	var err error
	this.socket, err = net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Failed to start listener")
		return nil, false
	}

	go this.run()

	return this, true
}

func (this *Listener) run() {
	for {
		s, err := this.socket.Accept()

		// TODO: differentiate between errors (ie: break for closure of listening socket)
		if err != nil {
			fmt.Println("Error with listener")
			return
		}

		_ = NewPeer(s, this.context)

	}
}

func (this *Network) SpawnListener(port string) (success bool) {
	listener, listening := NewListener(this.context, port)

	if listening {
		this.listeners = append(this.listeners, listener)
		return true
	}

	return false
}

func (this *Network) SpawnPeer(address string) (success bool) {

	socket, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println("Failed to connect to peer")
		return false
	}
	fmt.Println("Successfully connected to peer")

	_ = NewPeer(socket, this.context)

	return true
}
