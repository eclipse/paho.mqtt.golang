package paho

import (
	"sync"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

// Persistence is an interface of the functions for a struct
// that is used to persist ControlPackets.
// Open() is an initialiser to prepare the Persistence for use
// Put() takes a uint16 which is a messageid and a ControlPacket
// to persist against that messageid
// Get() takes a uint16 which is a messageid and returns the
// persisted ControlPacket from the Persistence for that messageid
// All() returns a slice of all ControlPackets persisted
// Delete() takes a uint16 which is a messageid and deletes the
// associated stored ControlPacket from the Persistence
// Close() closes the Persistence
// Reset() clears the Persistence and prepares it to be reused
type Persistence interface {
	Open()
	Put(uint16, p.ControlPacket)
	Get(uint16) p.ControlPacket
	All() []p.ControlPacket
	Delete(uint16)
	Close()
	Reset()
}

// MemoryPersistence is an implementation of a Persistence
// that stores the ControlPackets in memory using a map
type MemoryPersistence struct {
	sync.RWMutex
	packets map[uint16]p.ControlPacket
}

// Open is the library provided MemoryPersistence's implementation of
// the required interface function()
func (m *MemoryPersistence) Open() {
	m.Lock()
	m.packets = make(map[uint16]p.ControlPacket)
	m.Unlock()
}

// Put is the library provided MemoryPersistence's implementation of
// the required interface function()
func (m *MemoryPersistence) Put(id uint16, cp p.ControlPacket) {
	m.Lock()
	m.packets[id] = cp
	m.Unlock()
}

// Get is the library provided MemoryPersistence's implementation of
// the required interface function()
func (m *MemoryPersistence) Get(id uint16) p.ControlPacket {
	m.RLock()
	defer m.RUnlock()
	return m.packets[id]
}

// All is the library provided MemoryPersistence's implementation of
// the required interface function()
func (m *MemoryPersistence) All() []p.ControlPacket {
	m.Lock()
	defer m.RUnlock()
	ret := make([]p.ControlPacket, len(m.packets))

	for _, cp := range m.packets {
		ret = append(ret, cp)
	}

	return ret
}

// Delete is the library provided MemoryPersistence's implementation of
// the required interface function()
func (m *MemoryPersistence) Delete(id uint16) {
	m.Lock()
	delete(m.packets, id)
	m.Unlock()
}

// Close is the library provided MemoryPersistence's implementation of
// the required interface function()
func (m *MemoryPersistence) Close() {
	m.Lock()
	m.packets = nil
	m.Unlock()
}

// Reset is the library provided MemoryPersistence's implementation of
// the required interface function()
func (m *MemoryPersistence) Reset() {
	m.Lock()
	m.packets = make(map[uint16]p.ControlPacket)
	m.Unlock()
}
