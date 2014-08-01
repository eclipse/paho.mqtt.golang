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
	"encoding/binary"
	"time"
)

type MsgType byte
type QoS byte
type ConnRC int8

type Message struct {
	lastActivity time.Time
	header       uint8
	remlen       uint32
	messageId    MId
	vheader      []byte
	payload      []byte
}

type sendable struct {
	m *Message
	r chan Receipt
}

// Create a default PUBLISH Message with the specified payload
// If message == nil, create a zero length message
// Defaults: QoS=1, Retained=False
func NewMessage(message []byte) *Message {
	m := newMsg(PUBLISH, false, QOS_ONE, false)
	if message == nil {
		m.appendPayloadField([]byte{})
	} else {
		m.appendPayloadField(message)
	}
	return m
}

//newMsg takes a message type, a boolean indicating if this message
//is a duplicate, a QoS value for the message and a boolean indicating
//if this message is to be retained by the server. It returns a
//pointer to a Message initialized with these values
func newMsg(msgtype MsgType, duplicate bool, qos QoS, retained bool) *Message {
	m := &Message{}
	m.setMsgType(msgtype)
	m.setDupFlag(duplicate)
	m.SetQoS(qos)
	m.SetRetainedFlag(retained)
	return m
}

//newConnectMsg takes a boolean indicating whether this connection should
//be a clean session, a boolean indicating the existence of a "will"
//message, a QoS value for the will message, a boolean indicating if the
//will message should be retained, the topic the will message will be
//published on, the content of the will message as a slice of bytes,
//the clientid to be used for the connection, username and password
//strings and a keep alive value.
//It returns a pointer to a Message initialized with these values

func newConnectMsg(
	cleanSession bool,
	will bool,
	willqos QoS,
	willretained bool,
	willtopic string,
	willmessage []byte,
	clientid,
	user,
	password string,
	keepalive uint16) *Message {
	m := newMsg(CONNECT, false, 0, false)

	m.remlen = uint32(0)

	/* Protocol Name */
	m.appendVHeaderField("MQIsdp")

	/* Protocol Version */
	m.vheader = append(m.vheader, 0x03)

	/* Connect Byte */
	b := byte(0)
	if cleanSession {
		b |= 0x02
	}
	if will {
		b |= 0x04
	}
	b |= byte(willqos) << 3
	if willretained {
		b |= 0x20
	}

	m.appendPayloadSizedField(clientid)
	if will {
		m.appendPayloadSizedField(willtopic)
		m.appendPayloadSizedField(willmessage)
	}

	if user != "" {
		b |= 0x80
		m.appendPayloadSizedField(user)
		//mustn't have password without user as well
		if password != "" {
			b |= 0x40
			m.appendPayloadSizedField(password)
		}
	}

	m.vheader = append(m.vheader, b)

	/* Keep Alive Time Interval (2 Bytes [MSB, LSB]) */
	m.appendVHeaderField(keepalive)

	numbytes := uint(len(m.vheader) + len(m.payload))
	m.remlen = encodeLength(numbytes)
	return m
}

func newConnectMsgFromOptions(options ClientOptions) *Message {
	m := newMsg(CONNECT, false, 0, false)

	m.remlen = uint32(0)

	/* Protocol Name */
	m.appendVHeaderField("MQIsdp")

	/* Protocol Version */
	m.vheader = append(m.vheader, 0x03)

	/* Connect Byte */
	b := byte(0)
	if options.cleanSession {
		b |= 0x02
	}
	if options.willEnabled {
		b |= 0x04
	}
	b |= byte(options.willQos) << 3
	if options.willRetained {
		b |= 0x20
	}

	m.appendPayloadSizedField(options.clientId)
	if options.willEnabled {
		m.appendPayloadSizedField(options.willTopic)
		m.appendPayloadSizedField(options.willPayload)
	}

	if options.username != "" {
		b |= 0x80
		m.appendPayloadSizedField(options.username)
		//mustn't have password without user as well
		if options.password != "" {
			b |= 0x40
			m.appendPayloadSizedField(options.password)
		}
	}

	m.vheader = append(m.vheader, b)

	/* Keep Alive Time Interval (2 Bytes [MSB, LSB]) */
	m.appendVHeaderField(uint16(options.keepAlive))

	numbytes := uint(len(m.vheader) + len(m.payload))
	m.remlen = encodeLength(numbytes)
	return m
}

//newDisconnectMsg returns a pointer to a DISCONNECT message
func newDisconnectMsg() *Message {
	m := newMsg(DISCONNECT, false, QOS_ZERO, false)
	m.remlen = uint32(0)
	return m
}

//encodeUint16 takes a uint16 and returns a slice of bytes
//representing its value in network order
func encodeUint16(num uint16) []byte {
	encNum := make([]byte, 2)
	binary.BigEndian.PutUint16(encNum, num)
	return encNum
}

//setTime sets the lastActivity field of a Message to time.Now()
func (m *Message) setTime() {
	m.lastActivity = time.Now()
}

//appendVHeaderField takes a single value of varying types and
//encodes the appropriate value in the variable header field
//of the Message
func (m *Message) appendVHeaderField(field interface{}) {
	switch f := field.(type) {
	case MId:
		m.vheader = append(m.vheader, encodeUint16(uint16(f))...)
	case uint16:
		m.vheader = append(m.vheader, encodeUint16(f)...)
	case string:
		m.vheader = append(m.vheader, encodeUint16(uint16(len(f)))...)
		m.vheader = append(m.vheader, []byte(f)...)
	case []byte:
		m.vheader = append(m.vheader, encodeUint16(uint16(len(f)))...)
		m.vheader = append(m.vheader, f...)
	}
}

//appendPayloadSizedField takes a single value of string or
//slice of bytes and appends a length prefixed value into the payload
//field of the Message
func (m *Message) appendPayloadSizedField(field interface{}) {
	switch f := field.(type) {
	case string:
		m.payload = append(m.payload, encodeUint16(uint16(len(f)))...)
		m.payload = append(m.payload, []byte(f)...)
	case []byte:
		m.payload = append(m.payload, encodeUint16(uint16(len(f)))...)
		m.payload = append(m.payload, f...)
	}
}

//appendPayloadField takes a single value of string or slice
//of bytes and appends the value to the payload field of the Message
func (m *Message) appendPayloadField(field interface{}) {
	switch f := field.(type) {
	case string:
		m.payload = append(m.payload, []byte(f)...)
	case []byte:
		m.payload = append(m.payload, f...)
	}
}

//Topic returns the topic of the Message as encoded in the variable header
func (m *Message) Topic() string {
	// skip 2 topic length bytes
	if m.QoS() == QOS_ZERO {
		return string(m.vheader[2:])
	}
	// do not include message id bytes in qos 1 and 2
	// do not include topic length bytes
	return string(m.vheader[2:len(m.vheader)])
}

//Payload returns a slice of bytes containing the payload of the Message
func (m *Message) Payload() []byte {
	return m.payload
}

//timeout returns a uint16 indicating the timeout value of the Message
//as encoded in the variable header.
func (m *Message) timeout() uint16 {
	timeout := uint16(0)
	timeout |= uint16(m.vheader[10]) << 8
	timeout |= uint16(m.vheader[11])
	return timeout
}

//connRC returns the return code from a CONNACK Message as encoded
//in the variable header
func (m *Message) connRC() ConnRC {
	rc := ConnRC(m.vheader[1])
	return rc
}

//msgType returns the message type of the Message as encoded in the
//fixed header
func (m *Message) msgType() MsgType {
	return MsgType((m.header & 0xF0) >> 4)
}

//setRemLen takes a uint32 indicating the remaining length of the
//Message after the fixed header and encodes this value
func (m *Message) setRemLen(length uint32) {
	m.remlen = length
}

//setMsgType takes a message type and sets the receiving Message
//to this type
func (m *Message) setMsgType(msgtype MsgType) {
	m.header &= 0x0F
	m.header |= (uint8(msgtype) << 4)
}

//DupFlag returns the boolean value of the duplicate message flag as encoded
//in the fixed header
func (m *Message) DupFlag() bool {
	return (m.header & 0x08) == 0x08
}

//setDupFlag takes a boolean value indicating whether this Message
//is a duplicated message
func (m *Message) setDupFlag(isDup bool) {
	m.header &= 0xF7
	if isDup {
		m.header |= 0x08
	}
}

//QoS returns the QoS value of the Message as encoded in the fixed header
func (m *Message) QoS() QoS {
	return QoS((m.header & 0x06) >> 1)
}

//setQoS takes a QoS value and encodes this value in the fixed header of
//the Message
func (m *Message) SetQoS(qos QoS) {
	m.header &= 0xF9
	m.header |= (uint8(qos) << 1)
}

//RetainedFlag returns a boolean value indicating whether this message
//was a retained message
func (m *Message) RetainedFlag() bool {
	return (m.header & 0x01) == 0x01
}

//SetRetainedFlag takes a boolean value indicating whether the server
//should retain the message and encodes it in the fixed header
func (m *Message) SetRetainedFlag(isRetained bool) {
	m.header &= 0xFE
	if isRetained {
		m.header |= 0x01
	}
}

//remLen returns a uint of the remaining length value encoded in the Message
func (m *Message) remLen() uint {
	return decodeRemLen(m.remlen)
}

//MsgId returns a MId containing the message id of the Message
func (m *Message) MsgId() MId {
	return m.messageId
}

//setMsgId takes a MId and sets the message id field of the Message to this
//value
func (m *Message) setMsgId(id MId) {
	if m.QoS() != QOS_ZERO || m.msgType() != PUBLISH {
		m.messageId = id
	}
}

//newSubscribeMsg takes a list of TopicFilter
//and returns a Message initialized with these values
func newSubscribeMsg(filters ...*TopicFilter) *Message {

	m := newMsg(SUBSCRIBE, false, QOS_ONE, false)

	for i := range filters {
		m.appendPayloadSizedField(filters[i].string)
		m.appendPayloadField([]byte{byte(filters[i].QoS)})
	}

	numbytes := uint(len(m.vheader) + len(m.payload) + 2)
	m.remlen = encodeLength(numbytes)
	return m
}

//newUnsubscribeMsg takes a list of strings of topics that are to be
//unsubscribed from and returns a pointer to a Message initialized with
//these values
func newUnsubscribeMsg(subscriptions ...string) *Message {
	m := newMsg(UNSUBSCRIBE, false, QOS_ONE, false)

	for i := range subscriptions {
		m.appendPayloadSizedField(subscriptions[i])
	}

	numbytes := uint(len(m.vheader) + len(m.payload) + 2)
	m.remlen = encodeLength(numbytes)
	return m
}

//newPublishMsg takes a QoS value, a string of the topic the message
//is to be published on and a slice of bytes of the payload of the message
//It returns a pointer to a Message initialized with these values
func newPublishMsg(qos QoS, topic string, payload []byte) *Message {
	m := newMsg(PUBLISH, false, qos, false)

	/* vheader, topic */
	m.appendVHeaderField(topic)
	/* vheader, message id if QOS > 0 */

	m.appendPayloadField(payload)

	numbytes := uint(len(m.vheader) + len(m.payload))
	if qos > QOS_ZERO {
		numbytes += 2
	}
	m.remlen = encodeLength(numbytes)
	return m
}

//newPubRelMsg returns a pointer to a new PUBREL Message
func newPubRelMsg() *Message {
	m := newMsg(PUBREL, false, QOS_ONE, false)
	m.remlen = encodeLength(uint(len(m.vheader)) + 2)
	return m
}

//newPubRedMsg returns a pointer to a new PUBREC Message
func newPubRecMsg() *Message {
	m := newMsg(PUBREC, false, QOS_ZERO, false)
	m.remlen = encodeLength(uint(len(m.vheader)) + 2)
	return m
}

//newPubCompMsg returns a pointer to a new PUBCOMP Message
func newPubCompMsg() *Message {
	m := newMsg(PUBCOMP, false, QOS_ZERO, false)
	m.remlen = encodeLength(uint(len(m.vheader)) + 2)
	return m
}

//newPubAckMsg returns a pointer to a new PUBACK Message
func newPubAckMsg() *Message {
	m := newMsg(PUBACK, false, QOS_ZERO, false)
	m.remlen = encodeLength(uint(len(m.vheader)) + 2)
	return m
}

//newSubAckMsg returns a pointer to a new SUBACK Message
func newSubackMsg() *Message {
	m := newMsg(SUBACK, false, QOS_ZERO, false)
	m.remlen = encodeLength(uint(len(m.vheader)) + 2)
	return m
}

//newUnsubackMsg returns a pointer to a new UNSUBACK Message
func newUnsubackMsg() *Message {
	m := newMsg(UNSUBACK, false, QOS_ZERO, false)
	m.remlen = encodeLength(uint(len(m.vheader)) + 2)
	return m
}

//encodeLength takes a uint value and returns uint32 of that value
//in the encoding mechanism defined for MQTT remaining lengths
func encodeLength(length uint) uint32 {
	value := uint32(0)
	digit := uint32(0)
	X := uint32(length)
	for X > 0 {
		if value != 0 {
			value <<= 8
		}
		digit = X % 128
		X /= 128
		if X > 0 {
			digit |= 0x80
		}
		value |= uint32(digit)
	}
	return value
}

//decodeRemLen takes a uint32 of the remaining length encoded as
//defined by MQTT and returns a uint of the value
func decodeRemLen(remlen uint32) uint {
	multiplier := uint32(1)

	value := uint32(uint8(remlen>>(8*3))&0x7F) * multiplier
	if remlen>>(8*3) != 0 {
		multiplier *= 128
	}

	value += uint32(uint8(remlen>>(8*2))&0x7F) * multiplier
	if remlen>>(8*2) != 0 {
		multiplier *= 128
	}

	value += uint32(uint8(remlen>>(8*1))&0x7F) * multiplier
	if remlen>>(8*1) != 0 {
		multiplier *= 128
	}

	value += uint32(uint8(remlen>>(8*0))&0x7F) * multiplier

	return uint(value)
}

/* MsgType */
const (
	/* 0x00 is reserved */
	CONNECT     MsgType = 0x01
	CONNACK     MsgType = 0x02
	PUBLISH     MsgType = 0x03
	PUBACK      MsgType = 0x04
	PUBREC      MsgType = 0x05
	PUBREL      MsgType = 0x06
	PUBCOMP     MsgType = 0x07
	SUBSCRIBE   MsgType = 0x08
	SUBACK      MsgType = 0x09
	UNSUBSCRIBE MsgType = 0x0A
	UNSUBACK    MsgType = 0x0B
	PINGREQ     MsgType = 0x0C
	PINGRESP    MsgType = 0x0D
	DISCONNECT  MsgType = 0x0E
	/* 0x0F is reserved */
)

/* QoS LEVEL */
const (
	QOS_ZERO QoS = 0
	QOS_ONE  QoS = 1
	QOS_TWO  QoS = 2
)

/* Connection Return Codes */
const (
	CONN_FAILURE           ConnRC = -1
	CONN_ACCEPTED          ConnRC = 0x00
	CONN_REF_BAD_PROTO_VER ConnRC = 0x01
	CONN_REF_ID_REJ        ConnRC = 0x02
	CONN_REF_SERV_UNAVAIL  ConnRC = 0x03
	CONN_REF_BAD_USER_PASS ConnRC = 0x04
	CONN_REF_NOT_AUTH      ConnRC = 0x05
)
