package packets

import (
	"bufio"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiveConnection(t *testing.T) {
	x := NewConnect(
		ClientID("testClient"),
		Username("testUser"),
		KeepAlive(30),
	)

	sExpiryInterval := uint32(30)
	x.Properties.SessionExpiryInterval = &sExpiryInterval

	conn, err := net.Dial("tcp", "127.0.0.1:1883")
	require.Nil(t, err)

	err = x.Send(conn)
	require.Nil(t, err)

	p, err := ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, CONNACK, p.Type)

	s := NewSubscribe(
		Sub("test", SubOptions{QoS: 1}),
	)
	s.PacketID = 1

	err = s.Send(conn)
	require.Nil(t, err)

	sa, err := ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, SUBACK, sa.Type)

	pb := NewPublish(
		Message("test", 0, false, []byte("Test message")),
	)
	pb.PacketID = 2

	err = pb.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBLISH, p.Type)
	fmt.Println(string(p.Content.(*Publish).Payload))

	pb = NewPublish(
		Message("testqos1", 1, false, []byte("Test message")),
	)

	pb.PacketID = 3

	err = pb.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBACK, p.Type)

	pb = NewPublish(
		Message("testqos2", 2, false, []byte("Test message")),
	)
	pb.PacketID = 4

	err = pb.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBREC, p.Type)

	pr := NewControlPacket(PUBREL)
	pr.Content.(*Pubrel).PacketID = 4

	err = pr.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBCOMP, p.Type)
}
