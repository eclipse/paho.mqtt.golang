package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

type noopPersistence struct{}

func (n *noopPersistence) Open() {}

func (n *noopPersistence) Put(id uint16, cp packets.ControlPacket) {}

func (n *noopPersistence) Get(id uint16) packets.ControlPacket {
	return packets.ControlPacket{}
}

func (n *noopPersistence) All() []packets.ControlPacket {
	return nil
}

func (n *noopPersistence) Delete(id uint16) {}

func (n *noopPersistence) Close() {}

func (n *noopPersistence) Reset() {}
