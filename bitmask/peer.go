package bitmask

import (
	"fmt"
	"math"
	"net"
	"unsafe"
)

type Peer struct {
	socket  net.Conn
	context *Context

	Version   Version
	Mode      Mode
	SubnetMap map[Subnet]bool
	UserAgent string

	AdListener      *BufferListener
	PermReqListener *BufferListener
}

func NewPeer(socket net.Conn, context *Context) (this *Peer) {
	this = new(Peer)
	this.socket = socket
	this.context = context

	this.SubnetMap = make(map[Subnet]bool)

	this.AdListener = context.AdPackets.SpawnListener()
	this.PermReqListener = context.PermReqPackets.SpawnListener()

	go this.run()

	return this
}

func (this *Peer) Send(message []byte) (err error) {
	_, err = this.socket.Write(message)

	return err
}

func (this *Peer) run() {
	defer this.socket.Close()
	defer this.cleanup()

	var header [12]byte // command + len(payload) + checksum
	var commandType, payloadLen, checkSum uint32
	for {
		length, err := this.socket.Read(header[:])

		if length != 12 {
			fmt.Println("Invalid header => less than 12 bytes")
			return
		}

		if err != nil {
			fmt.Println("Peer disconnected before/during new command")
			return
		}

		commandType = *(*uint32)(unsafe.Pointer(&header[0]))
		payloadLen = *(*uint32)(unsafe.Pointer(&header[4]))
		checkSum = *(*uint32)(unsafe.Pointer(&header[8]))

		command := Command(commandType)

		payload := make([]byte, payloadLen) // TODO: make this more efficient
		length, err = this.socket.Read(payload[:])

		if uint32(length) != payloadLen {
			fmt.Println("Invalid header => payload length mismatch")
			return
		}
		if err != nil {
			fmt.Println("Payload read error")
			return
		}

		// TODO: Verify payload with checksum

		var canParse bool

		switch command {
		case C_VERSION:
			canParse = this.parseVersion(payload)
		default:
			fmt.Println("Invalid header => non-existant command")
			return
		}

		if !canParse {
			fmt.Println("Failed to parse")
			return
		}
	}
}

func (this *Peer) cleanup() {
	this.AdListener.Stop()
	this.PermReqListener.Stop()
}

func (this *Peer) parseVersion(payload []byte) (canParse bool) {
	length := len(payload)

	if length < 8 {
		fmt.Println("VERSION payload too short -> < 8 bytes")
		return false
	}

	version := float64(*(*uint32)(unsafe.Pointer(&payload[0])))
	this.Version = Version(uint32(math.Min(version, float64(this.context.Version))))

	this.Mode = Mode(*(*uint32)(unsafe.Pointer(&payload[4])))

	start := 8
	subnetCount, start := parseVarInt(payload, start)
	if start == -1 {
		fmt.Println("Error while parsing VERSION subnet count")
		return false
	}

	// Assumption: subnetCount is never > 2^31 - 1 (ie: 32-bit systems won't overflow to negatives)

	for i := 0; i < int(subnetCount); i++ {

		if start+4 > length {
			fmt.Println("VERSION payload too short -> not enough subnets")
			return false
		}

		this.SubnetMap[Subnet(*(*uint32)(unsafe.Pointer(&payload[start])))] = true

		start += 4
	}

	str, _ := parseVarStr(payload, start)
	if str == nil {
		fmt.Println("Error while parsing VERSION UserAgent")
		return false
	}
	this.UserAgent = *str

	return true
}
