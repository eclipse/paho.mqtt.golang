package paho

import p "github.com/eclipse/paho.mqtt.golang/packets"

type noopPersistence struct{}

func (n *noopPersistence) Open() {}

func (n *noopPersistence) Put(id uint16, cp p.ControlPacket) {}

func (n *noopPersistence) Get(id uint16) p.ControlPacket {
	return p.ControlPacket{}
}

func (n *noopPersistence) All() []p.ControlPacket {
	return nil
}

func (n *noopPersistence) Delete(id uint16) {}

func (n *noopPersistence) Close() {}

func (n *noopPersistence) Reset() {}
