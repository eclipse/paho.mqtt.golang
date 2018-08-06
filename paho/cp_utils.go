package paho

// Byte is a helper function that take a byte and returns
// a pointer to a byte of that value
func Byte(b byte) *byte {
	return &b
}

// Uint32 is a helper function that take a uint32 and returns
// a pointer to a uint32 of that value
func Uint32(u uint32) *uint32 {
	return &u
}

// Uint16 is a helper function that take a uint16 and returns
// a pointer to a uint16 of that value
func Uint16(u uint16) *uint16 {
	return &u
}

// BoolToByte is a helper function that take a bool and returns
// a pointer to a byte of value 1 if true or 0 if false
func BoolToByte(b bool) *byte {
	var v byte
	if b {
		v = 1
	}
	return &v
}
