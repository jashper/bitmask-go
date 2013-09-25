package bitmask

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Peer struct {
	socket  net.Conn
	context *Context
}

func NewPeer(socket net.Conn, context *Context) *Peer {
	this := new(Peer)
	this.socket = socket
	this.context = context
	go this.run()

	return this
}

func (this *Peer) Send(message []byte) error {
	_, err := this.socket.Write(message)

	return err
}

func (this *Peer) run() {
	defer this.socket.Close()
	fmt.Println("New user connected")

	var header [12]byte // command + len(payload) + checksum
	var commandType, payloadLen, checkSum uint32
	for {
		_, err := this.socket.Read(header[:])

		if err != nil {
			fmt.Println("User disconnected")
			return
		}

		buffer := bytes.NewBuffer(header[0:4])
		binary.Read(buffer, binary.BigEndian, &commandType)
		buffer = bytes.NewBuffer(header[4:8])
		binary.Read(buffer, binary.BigEndian, &payloadLen)
		buffer = bytes.NewBuffer(header[8:12])
		binary.Read(buffer, binary.BigEndian, &checkSum)

		command := Command(commandType)

		payload := make([]byte, payloadLen)
		_, err = this.socket.Read(payload[:])

		// TODO: Verify payload with checksum

		switch command {
		case VERSION:
			this.parseVersion(payload)
		default:
			fmt.Println("Invalid Command")
		}
	}
}

func (this *Peer) parseVersion(payload []byte) {

}
