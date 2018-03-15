package paho

import (
	"context"
	"sync"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

const (
	midMin uint16 = 1
	midMax uint16 = 65535
)

type MIDService interface {
	Request(*CPContext) uint16
	Get(uint16) *CPContext
	Free(uint16)
	Clear()
}

type CPContext struct {
	Return  chan p.ControlPacket
	Context context.Context
}

type MIDs struct {
	sync.Mutex
	index map[uint16]*CPContext
}

func (m *MIDs) Request(c *CPContext) uint16 {
	m.Lock()
	defer m.Unlock()
	for i := midMin; i < midMax; i++ {
		if _, ok := m.index[i]; !ok {
			m.index[i] = c
			return i
		}
	}
	return 0
}

func (m *MIDs) Get(i uint16) *CPContext {
	m.Lock()
	defer m.Unlock()
	return m.index[i]
}

func (m *MIDs) Free(i uint16) {
	m.Lock()
	delete(m.index, i)
	m.Unlock()
}

func (m *MIDs) Clear() {
	m.index = make(map[uint16]*CPContext)
}
