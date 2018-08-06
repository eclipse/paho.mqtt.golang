package packets

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeVBI127(t *testing.T) {
	b := encodeVBI(127)

	require.Len(t, b, 1)
	assert.Equal(t, byte(127), b[0])
}

func TestEncodeVBI128(t *testing.T) {
	b := encodeVBI(128)

	require.Len(t, b, 2)
	assert.Equal(t, byte(0x80), b[0])
	assert.Equal(t, byte(0x01), b[1])
}

func TestEncodeVBI16383(t *testing.T) {
	b := encodeVBI(16383)

	require.Len(t, b, 2)
	assert.Equal(t, byte(0xff), b[0])
	assert.Equal(t, byte(0x7f), b[1])
}

func TestEncodeVBI16384(t *testing.T) {
	b := encodeVBI(16384)

	require.Len(t, b, 3)
	assert.Equal(t, byte(0x80), b[0])
	assert.Equal(t, byte(0x80), b[1])
	assert.Equal(t, byte(0x01), b[2])
}

func TestEncodeVBI2097151(t *testing.T) {
	b := encodeVBI(2097151)

	require.Len(t, b, 3)
	assert.Equal(t, byte(0xff), b[0])
	assert.Equal(t, byte(0xff), b[1])
	assert.Equal(t, byte(0x7f), b[2])
}

func TestEncodeVBI2097152(t *testing.T) {
	b := encodeVBI(2097152)

	require.Len(t, b, 4)
	assert.Equal(t, byte(0x80), b[0])
	assert.Equal(t, byte(0x80), b[1])
	assert.Equal(t, byte(0x80), b[2])
	assert.Equal(t, byte(0x01), b[3])
}

func TestEncodeVBIMax(t *testing.T) {
	b := encodeVBI(268435455)

	require.Len(t, b, 4)
	assert.Equal(t, byte(0xff), b[0])
	assert.Equal(t, byte(0xff), b[1])
	assert.Equal(t, byte(0xff), b[2])
	assert.Equal(t, byte(0x7f), b[3])
}

func TestDecodeVBI12(t *testing.T) {
	x, err := decodeVBI(bytes.NewBuffer([]byte{0x0C}))

	require.Nil(t, err)
	assert.Equal(t, 12, x)
}

func TestDecodeVBI127(t *testing.T) {
	x, err := decodeVBI(bytes.NewBuffer([]byte{0xff}))

	require.Nil(t, err)
	assert.Equal(t, 127, x)
}
func TestDecodeVBI128(t *testing.T) {
	x, err := decodeVBI(bytes.NewBuffer([]byte{0x80, 0x01}))

	require.Nil(t, err)
	assert.Equal(t, 128, x)
}
func TestDecodeVBI16384(t *testing.T) {
	x, err := decodeVBI(bytes.NewBuffer([]byte{0x80, 0x80, 0x01}))

	require.Nil(t, err)
	assert.Equal(t, 16384, x)
}
func TestDecodeVBIMax(t *testing.T) {
	x, err := decodeVBI(bytes.NewBuffer([]byte{0xff, 0xff, 0xff, 0x7f}))

	require.Nil(t, err)
	assert.Equal(t, 268435455, x)
}

func TestNewControlPacketConnect(t *testing.T) {
	var b bytes.Buffer
	x := NewControlPacket(CONNECT)

	require.Equal(t, CONNECT, x.Type)

	x.Content.(*Connect).KeepAlive = 30
	x.Content.(*Connect).ClientID = "testClient"
	x.Content.(*Connect).UsernameFlag = true
	x.Content.(*Connect).Username = "testUser"
	sExpiryInterval := uint32(30)
	x.Content.(*Connect).Properties.SessionExpiryInterval = &sExpiryInterval

	_, err := x.WriteTo(&b)

	require.Nil(t, err)
	assert.Len(t, b.Bytes(), 40)
}

func TestReadPacketConnect(t *testing.T) {
	p := []byte{16, 38, 0, 4, 77, 81, 84, 84, 5, 128, 0, 30, 5, 17, 0, 0, 0, 30, 0, 10, 116, 101, 115, 116, 67, 108, 105, 101, 110, 116, 0, 8, 116, 101, 115, 116, 85, 115, 101, 114}

	c, err := ReadPacket(bufio.NewReader(bytes.NewReader(p)))

	require.Nil(t, err)
	assert.Equal(t, uint16(30), c.Content.(*Connect).KeepAlive)
	assert.Equal(t, "testClient", c.Content.(*Connect).ClientID)
	assert.Equal(t, true, c.Content.(*Connect).UsernameFlag)
	assert.Equal(t, "testUser", c.Content.(*Connect).Username)
	assert.Equal(t, uint32(30), *c.Content.(*Connect).Properties.SessionExpiryInterval)
}

func TestReadStringWriteString(t *testing.T) {
	var b bytes.Buffer
	writeString("Test string", &b)

	s, err := readString(&b)
	require.Nil(t, err)
	assert.Equal(t, "Test string", s)
}
