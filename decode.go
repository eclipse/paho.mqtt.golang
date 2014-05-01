/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"bufio"
	"encoding/binary"
)

//decodeQos returns the QoS of the message
func decodeQos(header byte) QoS {
	qos := (header & 0x06) >> 1
	return QoS(qos)
}

//decodeMsgType returns the type of the message
func decodeMsgType(header byte) MsgType {
	mtype := (header & 0xF0) >> 4
	return MsgType(mtype)
}

// return (number of bytes needed, the remaining length)
func decodeRemlen(bytes []byte) (int, uint32) {
	// bytes[1,2,3,4] are the only ones we could care about
	idx := uint32(0)
	mult := uint32(1)
	value := uint32(0)
	for {
		idx++
		digit := uint32((bytes[idx]))
		value += (digit & 127) * mult
		mult *= 128
		if digit&128 == 0 {
			break
		}
	}
	return int(idx), value
}

func decodeRemlenFromNetwork(src *bufio.ReadWriter) ([]byte, uint32) {
	var bytes []byte
	var rLength uint32
	var count int
	var multiplier uint32 = 1
	var digit byte
	count = 1
	for {
		digit, _ = src.ReadByte()
		bytes = append(bytes, digit)
		rLength += uint32(digit&127) * multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier *= 128
		count++
	}
	return bytes, rLength
}

// return length of topic string, the topic string
func decodeTopic(bytes []byte) (tlen uint16, t string) {
	tlen = binary.BigEndian.Uint16(bytes[:2])
	t = string(bytes[2 : 2+tlen])
	return tlen, t
}

//decode takes a slice of bytes as received over the network
//and returns a Message pointer to a Message struct
func decode(bytes []byte) *Message {
	m := &Message{}

	m.SetQoS(decodeQos(bytes[0]))

	m.setMsgType(decodeMsgType(bytes[0]))

	n, r := decodeRemlen(bytes)
	m.setRemLen(r)

	bytes = bytes[n+1:] // skip past fixed header and variable length byte(s)

	switch m.msgType() {
	case CONNACK:
		m.vheader = append(m.vheader, 0x00) // bytes[0] of vheader not used
		m.vheader = append(m.vheader, bytes[1])
		/* No Payload */

	case PINGRESP:
		/* No vheader */
		/* No Payload */

	case PUBACK:
		m.setMsgId(MId(binary.BigEndian.Uint16(bytes[:2])))
		/* No Payload */

	case PUBREC:
		m.setMsgId(MId(binary.BigEndian.Uint16(bytes[:2])))
		/* No Payload */

	case PUBREL:
		m.setMsgId(MId(binary.BigEndian.Uint16(bytes[:2])))
		/* No Payload */

	case PUBCOMP:
		m.setMsgId(MId(binary.BigEndian.Uint16(bytes[:2])))
		/* No Payload */

	case SUBACK:
		m.setMsgId(MId(binary.BigEndian.Uint16(bytes[:2])))
		m.appendPayloadField(bytes[2:])

	case PUBLISH:
		// we are past the fixed header and variable remlen

		// now bytes[0] and bytes[1] are the length of the topic name (n)
		// bytes[2]... bytes[n] are the topic string
		// bytes[n+1] and bytes[n+2] are message id IFF QoS > 0
		// bytes[n+3]+ are the payload (if any)

		// bytes 0 and 1 are topic length
		tlen, topic := decodeTopic(bytes)

		m.appendVHeaderField(topic) // auto inserts length

		bytes = bytes[tlen+2:]
		// bytes[0] is now either msgid OR payload

		// if QoS > 0, 2 message id bytes are after the topic string
		if m.QoS() != QOS_ZERO {
			m.setMsgId(MId(binary.BigEndian.Uint16(bytes[:2])))
			m.appendPayloadField(bytes[2:])
		} else {
			m.appendPayloadField(bytes[0:])
		}
	case UNSUBACK:
		m.setMsgId(MId(binary.BigEndian.Uint16(bytes[:2])))
		/* No Payload */
	}

	return m
}

//Bytes operates on a Message pointer and returns a slice of bytes
//representing the Message ready for transmission over the network
func (m *Message) Bytes() []byte {
	var b []byte
	b = append(b, m.header)
	if m.remlen&0xFF000000 > 0 {
		b = append(b, byte((m.remlen&0xFF000000)>>24))
		b = append(b, byte((m.remlen&0x00FF0000)>>16))
		b = append(b, byte((m.remlen&0x0000FF00)>>8))
		b = append(b, byte(m.remlen&0x000000FF))
	} else if m.remlen&0x00FF0000 > 0 {
		b = append(b, byte((m.remlen&0x00FF0000)>>16))
		b = append(b, byte((m.remlen&0x0000FF00)>>8))
		b = append(b, byte(m.remlen&0x000000FF))
	} else if m.remlen&0x0000FF00 > 0 {
		b = append(b, byte((m.remlen&0x0000FF00)>>8))
		b = append(b, byte(m.remlen&0x000000FF))
	} else {
		b = append(b, byte(m.remlen&0x000000FF))
	}
	for i := range m.vheader {
		b = append(b, m.vheader[i])
	}
	if m.MsgId() != 0 {
		mid := make([]byte, 2)
		binary.BigEndian.PutUint16(mid, uint16(m.messageId))
		b = append(b, mid...)
	}
	for i := range m.payload {
		b = append(b, m.payload[i])
	}

	return b
}
