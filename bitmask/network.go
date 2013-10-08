package bitmask

import (
	"fmt"
	"net"
)

type Network struct {
	context   *Context
	listeners []*NetListener
}

type NetListener struct {
	context *Context
	socket  net.Listener
}

func NewNetwork(context *Context) (this *Network) {
	this = new(Network)
	this.context = context
	this.listeners = make([]*NetListener, 0)

	this.SpawnListener(":7001")

	return this
}

func NewNetListener(context *Context, port string) (this *NetListener, success bool) {
	this = new(NetListener)
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

func (this *NetListener) run() {
	for {
		s, err := this.socket.Accept()

		// TODO: differentiate between errors (ie: break for closure of listening socket)
		if err != nil {
			fmt.Println("Error with listener")
			return
		}

		NewPeer(s, this.context)

	}
}

func (this *Network) SpawnListener(port string) (success bool) {
	listener, listening := NewNetListener(this.context, port)

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

	NewPeer(socket, this.context)

	return true
}
