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
	x := NewControlPacket(CONNECT)

	require.Equal(t, CONNECT, x.Type)

	x.Content.(*Connect).KeepAlive = 30
	x.Content.(*Connect).ClientID = "testClient"
	x.Content.(*Connect).UsernameFlag = true
	x.Content.(*Connect).Username = "testUser"
	sExpiryInterval := uint32(30)
	x.Content.(*Connect).IDVP.SessionExpiryInterval = &sExpiryInterval

	conn, err := net.Dial("tcp", "127.0.0.1:1883")
	require.Nil(t, err)

	err = x.Send(conn)
	require.Nil(t, err)

	p, err := ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, CONNACK, p.Type)

	s := NewControlPacket(SUBSCRIBE)
	s.Content.(*Subscribe).PacketID = 1
	s.Content.(*Subscribe).Subscriptions["test"] = 1

	err = s.Send(conn)
	require.Nil(t, err)

	sa, err := ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, SUBACK, sa.Type)

	pb := NewControlPacket(PUBLISH)
	pb.Payload = []byte("Test message")
	pb.Content.(*Publish).PacketID = 2
	pb.Content.(*Publish).Topic = "test"

	err = pb.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBLISH, p.Type)
	fmt.Println(string(p.Payload))

	pb = NewControlPacket(PUBLISH)
	pb.Payload = []byte("Test message")
	pb.Content.(*Publish).PacketID = 3
	pb.Content.(*Publish).Topic = "testqos1"
	pb.FixedHeader.Flags = 2

	err = pb.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBACK, p.Type)

	pb = NewControlPacket(PUBLISH)
	pb.Payload = []byte("Test message")
	pb.Content.(*Publish).PacketID = 4
	pb.Content.(*Publish).Topic = "testqos2"
	pb.FixedHeader.Flags = 4

	err = pb.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBREC, p.Type)

	pb = NewControlPacket(PUBREL)
	pb.Content.(*Pubrel).PacketID = 4

	err = pb.Send(conn)
	require.Nil(t, err)

	p, err = ReadPacket(bufio.NewReader(conn))
	require.Nil(t, err)
	assert.Equal(t, PUBCOMP, p.Type)
}
