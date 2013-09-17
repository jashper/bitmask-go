package node

import (
	"net"
)

type Connection struct {
	socket net.Conn
}

func New(socket net.Conn) (c *Connection) {
	c.socket = socket
	go c.start()
}

func (this *Connection) start() {

}
