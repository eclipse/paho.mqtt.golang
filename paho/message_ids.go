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

// MIDService defines the interface for a struct that handles the
// relationship between message ids and CPContexts
// Request() takes a *CPContext and returns a uint16 that is the
// messageid that should be used by the code that called Request()
// Get() takes a uint16 that is a messageid and returns the matching
// *CPContext that the MIDService has associated with that messageid
// Free() takes a uint16 that is a messageid and instructs the MIDService
// to mark that messageid as available for reuse
// Clear() should reset the internal state of the MIDService
type MIDService interface {
	Request(*CPContext) uint16
	Get(uint16) *CPContext
	Free(uint16)
	Clear()
}

// CPContext is the struct that is used to return responses to
// ControlPackets that have them, eg: the suback to a subscribe.
// The reponse packet is send down the Return channel and the
// Context is used to track timeouts.
type CPContext struct {
	Return  chan p.ControlPacket
	Context context.Context
}

// MIDs is the default MIDService provided by this library.
// It uses a map of uint16 to *CPContext to track responses
// to messages with a messageid
type MIDs struct {
	sync.Mutex
	index map[uint16]*CPContext
}

// Request is the library provided MIDService's implementation of
// the required interface function()
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

// Get is the library provided MIDService's implementation of
// the required interface function()
func (m *MIDs) Get(i uint16) *CPContext {
	m.Lock()
	defer m.Unlock()
	return m.index[i]
}

// Free is the library provided MIDService's implementation of
// the required interface function()
func (m *MIDs) Free(i uint16) {
	m.Lock()
	delete(m.index, i)
	m.Unlock()
}

// Clear is the library provided MIDService's implementation of
// the required interface function()
func (m *MIDs) Clear() {
	m.index = make(map[uint16]*CPContext)
}
