package paho

// Byte is a helper function that take a byte value and returns
func Byte(b byte) *byte {
	return &b
}

func Uint32(u uint32) *uint32 {
	return &u
}

func Uint16(u uint16) *uint16 {
	return &u
}

func BoolToByte(b bool) *byte {
	var v byte
	if b {
		v = 1
	}
	return &v
}
