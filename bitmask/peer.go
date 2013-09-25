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

func NewPeer(socket net.Conn, context *Context) (this *Peer) {
	this = new(Peer)
	this.socket = socket
	this.context = context

	go this.run()

	return this
}

func (this *Peer) Send(message []byte) (err error) {
	_, err = this.socket.Write(message)

	return err
}

func (this *Peer) run() {
	defer this.socket.Close()

	var header [12]byte // command + len(payload) + checksum
	var buffA, buffB, buffC bytes.Buffer
	var commandType, payloadLen, checkSum uint32
	for {
		length, err := this.socket.Read(header[:])

		if length != 12 {
			fmt.Println("Invalid header => less than 12 bytes")
			return
		}

		if err != nil {
			fmt.Println("Peer disconnected")
			return
		}

		buffA.Write(header[0:4])
		binary.Read(&buffA, binary.BigEndian, &commandType)
		buffA.Reset()

		buffB.Write(header[4:8])
		binary.Read(&buffB, binary.BigEndian, &payloadLen)
		buffB.Reset()

		buffC.Write(header[8:12])
		binary.Read(&buffC, binary.BigEndian, &checkSum)
		buffC.Reset()

		command := Command(commandType)

		payload := make([]byte, payloadLen) // TODO: make this more efficient
		length, err = this.socket.Read(payload[:])

		if uint32(length) != payloadLen {
			fmt.Println("Invalid header => payload length mismatch")
			return
		}

		// TODO: Verify payload with checksum

		switch command {
		case VERSION:
			this.parseVersion(payload)
		default:
			fmt.Println("Invalid header => non-existant command")
			return
		}
	}
}

func (this *Peer) parseVersion(payload []byte) {

}
