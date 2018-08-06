package packets

import (
	"bufio"
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

	// sExpiryInterval := uint32(30)
	// x.Properties.SessionExpiryInterval = &sExpiryInterval

	conn, err := net.Dial("tcp", "127.0.0.1:1883")
	require.Nil(t, err)

	_, err = x.WriteTo(conn)
	require.Nil(t, err)

	p, err := ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, CONNACK, p.Type)

	s := NewSubscribe(
		SubscribeSingle("test", SubOptions{QoS: 1}),
	)
	s.PacketID = 1

	_, err = s.WriteTo(conn)
	require.Nil(t, err)

	sa, err := ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, SUBACK, sa.Type)

	pb := NewPublish(
		Message("test", 0, false, []byte("Test message")),
	)
	pb.PacketID = 2

	_, err = pb.WriteTo(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBLISH, p.Type)

	pb = NewPublish(
		Message("testqos1", 1, false, []byte("Test message")),
	)

	pb.PacketID = 3

	_, err = pb.WriteTo(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBACK, p.Type)

	pb = NewPublish(
		Message("testqos2", 2, false, []byte("Test message")),
	)
	pb.PacketID = 4

	_, err = pb.WriteTo(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBREC, p.Type)

	pr := NewControlPacket(PUBREL)
	pr.Content.(*Pubrel).PacketID = 4

	_, err = pr.WriteTo(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBCOMP, p.Type)
}
