package node

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Manager struct {
	connections []*Connection
}

func New() (m *Manager) {
	m.connections = make([]*Connection)
}

func (this *Manager) AddConnection(c *Connection) {

}
