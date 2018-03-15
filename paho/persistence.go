package paho

import (
	"sync"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

type Persistence interface {
	Open()
	Put(uint16, p.ControlPacket)
	Get(uint16) p.ControlPacket
	All() []p.ControlPacket
	Delete(uint16)
	Close()
	Reset()
}

type MemoryPersistence struct {
	sync.RWMutex
	packets map[uint16]p.ControlPacket
}

func (m *MemoryPersistence) Open() {
	m.Lock()
	m.packets = make(map[uint16]p.ControlPacket)
	m.Unlock()
}

func (m *MemoryPersistence) Put(id uint16, cp p.ControlPacket) {
	m.Lock()
	m.packets[id] = cp
	m.Unlock()
}

func (m *MemoryPersistence) Get(id uint16) p.ControlPacket {
	m.RLock()
	defer m.RUnlock()
	return m.packets[id]
}

func (m *MemoryPersistence) All() []p.ControlPacket {
	m.Lock()
	defer m.RUnlock()
	ret := make([]p.ControlPacket, len(m.packets))

	for _, cp := range m.packets {
		ret = append(ret, cp)
	}

	return ret
}

func (m *MemoryPersistence) Delete(id uint16) {
	m.Lock()
	delete(m.packets, id)
	m.Unlock()
}

func (m *MemoryPersistence) Close() {
	m.Lock()
	m.packets = nil
	m.Unlock()
}

func (m *MemoryPersistence) Reset() {
	m.Lock()
	m.packets = make(map[uint16]p.ControlPacket)
	m.Unlock()
}
