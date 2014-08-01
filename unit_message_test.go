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

import "testing"

func Test_MessageTypeValues(t *testing.T) {
	if CONNECT != 0x1 {
		t.Errorf("CONNECT is %x instead of %x", CONNECT, 0x1)
	}
	if CONNACK != 0x2 {
		t.Errorf("CONNACK is %x instead of %x", CONNACK, 0x2)
	}
	if PUBLISH != 0x3 {
		t.Errorf("PUBLISH is %x instead of %x", PUBLISH, 0x3)
	}
	if PUBACK != 0x4 {
		t.Errorf("PUBACK is %x instead of %x", PUBACK, 0x4)
	}
	if PUBREC != 0x5 {
		t.Errorf("PUBREC is %x instead of %x", PUBREC, 0x5)
	}
	if PUBREL != 0x6 {
		t.Errorf("PUBREL is %x instead of %x", PUBREL, 0x6)
	}
	if PUBCOMP != 0x7 {
		t.Errorf("PUBCOMP is %x instead of %x", PUBCOMP, 0x7)
	}
	if SUBSCRIBE != 0x8 {
		t.Errorf("SUBSCRIBE is %x instead of %x", SUBSCRIBE, 0x8)
	}
	if SUBACK != 0x9 {
		t.Errorf("SUBACK is %x instead of %x", SUBACK, 0x9)
	}
	if UNSUBSCRIBE != 0xA {
		t.Errorf("UNSUBSCRIBE is %x instead of %x", UNSUBSCRIBE, 0xA)
	}
	if UNSUBACK != 0xB {
		t.Errorf("UNSUBACK is %x instead of %x", UNSUBACK, 0xB)
	}
	if PINGREQ != 0xC {
		t.Errorf("PINGREQ is %x instead of %x", PINGREQ, 0xC)
	}
	if PINGRESP != 0xD {
		t.Errorf("PINGRESP is %x instead of %x", PINGRESP, 0xD)
	}
	if DISCONNECT != 0xE {
		t.Errorf("DISCONNECT is %x instead of %x", DISCONNECT, 0xE)
	}
}

func Test_setMsgType(t *testing.T) {
	m := newMsg(CONNECT, false, QOS_ZERO, false)
	if m.header != 0x10 {
		t.Errorf("m.header has value %x instead of %x", m.header, 0x10)
	}
	m = newMsg(PINGREQ, false, QOS_ZERO, false)
	if m.header != 0xC0 {
		t.Errorf("m.header has value %x instead of %x", m.header, 0xC0)
	}
}

func Test_msgType(t *testing.T) {
	m := newMsg(PUBLISH, false, QOS_ZERO, false)
	if m.msgType() != PUBLISH {
		t.Errorf("m.msgType has value %x instead of %x", m.msgType(), PUBLISH)
	}
	m = newMsg(DISCONNECT, false, QOS_ZERO, false)
	if m.msgType() != DISCONNECT {
		t.Errorf("m.msgType has value %x instead of %x", m.msgType(), DISCONNECT)
	}
}

func Test_DupFlag(t *testing.T) {
	m := newMsg(PUBACK, false, QOS_ZERO, false)
	if m.header&0x08 != 0 {
		t.Errorf("False dup flag set wrong")
	}
	if m.DupFlag() {
		t.Errorf("False dup flag set wrong")
	}
	m = newMsg(PUBREC, true, QOS_ZERO, false)
	if (m.header & 0x08) != 0x08 {
		t.Errorf("True dup flag set wrong")
	}
	if !m.DupFlag() {
		t.Errorf("True dup flag set wrong")
	}
}

func Test_setDupFlag(t *testing.T) {
	m := newMsg(UNSUBSCRIBE, false, QOS_ZERO, false)
	m.setDupFlag(true)
	if m.header&0x08 != 0x08 {
		t.Errorf("setDupFlag true set wrong")
	}
	m.setDupFlag(false)
	if m.header&0x08 != 0 {
		t.Errorf("setDupFlag false set wrong")
	}

	m = newMsg(PINGRESP, true, QOS_ZERO, false)
	m.setDupFlag(false)
	if m.header&0x08 != 0 {
		t.Errorf("setDupFlag false set wrong")
	}
	m.setDupFlag(true)
	if m.header&0x08 != 0x08 {
		t.Errorf("setDupFlag true set wrong")
	}
}

func Test_QoS(t *testing.T) {
	m := newMsg(PUBREL, false, QOS_ZERO, false)
	if m.QoS() != QOS_ZERO {
		t.Errorf("QoS wrong value")
	}
	m = newMsg(PUBREL, false, QOS_ONE, false)
	if m.QoS() != QOS_ONE {
		t.Errorf("QoS wrong value")
	}
	m = newMsg(PUBREL, false, QOS_TWO, false)
	if m.QoS() != QOS_TWO {
		t.Errorf("QoS wrong value")
	}
}

func Test_setQoS(t *testing.T) {
	m := newMsg(UNSUBACK, false, QOS_ONE, false)
	m.SetQoS(QOS_ZERO)
	if m.QoS() != QOS_ZERO {
		t.Errorf("setQoS set wrong")
	}
	m.SetQoS(QOS_TWO)
	if m.QoS() != QOS_TWO {
		t.Errorf("setQoS set wrong")
	}
	if m.header != 0xB4 {
		t.Errorf("setQoS corrupted header, %x", m.header)
	}
}

func Test_RetainedFlag(t *testing.T) {
	m := newMsg(CONNACK, false, QOS_ZERO, false)
	if m.RetainedFlag() {
		t.Errorf("RetainedFlag got wrong value")
	}
	m = newMsg(CONNACK, false, QOS_ZERO, true)
	if !m.RetainedFlag() {
		t.Errorf("RetainedFlag got wrong value")
	}
}

func Test_setRetainedFlag(t *testing.T) {
	m := newMsg(CONNACK, false, QOS_ZERO, false)
	if m.RetainedFlag() {
		t.Errorf("RetainedFlag got wrong value")
	}
	m = newMsg(CONNACK, false, QOS_ZERO, true)
	if !m.RetainedFlag() {
		t.Errorf("RetainedFlag got wrong value")
	}
}

func Test_HeaderFuzz(t *testing.T) {
	m := newMsg(PUBCOMP, false, QOS_TWO, true)
	if m.header != 0x75 {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	m.SetQoS(QOS_ZERO)
	if m.header != 0x71 {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	m.SetRetainedFlag(false)
	if m.header != 0x70 {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	m.setDupFlag(true)
	if m.header != 0x78 {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	m.setMsgType(DISCONNECT)
	if m.header != 0xE8 {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	m.SetQoS(QOS_ONE)
	if m.header != 0xEA {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	m.SetRetainedFlag(true)
	if m.header != 0xEB {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	m.setMsgType(PUBREL)
	if m.header != 0x6B {
		t.Errorf("header fuzz test bad value: %x", m.header)
	}
	if m.msgType() != PUBREL {
		t.Errorf("msgType wrong type")
	}
	if m.DupFlag() != true {
		t.Errorf("DupFlag wrong value")
	}
	if m.QoS() != QOS_ONE {
		t.Errorf("QoS wrong value")
	}
	if m.RetainedFlag() != true {
		t.Errorf("RetainedFlag wrong value")
	}
}

func Test_encodeLength(t *testing.T) {
	m := map[uint]uint32{
		0:         0x00000000,
		1:         0x00000001,
		2:         0x00000002,
		3:         0x00000003,
		8:         0x00000008,
		15:        0x0000000F,
		16:        0x00000010,
		127:       0x0000007F,
		128:       0x00008001,
		16383:     0x0000FF7F,
		16384:     0x00808001,
		2097151:   0x00FFFF7F,
		2097152:   0x80808001,
		268435455: 0xFFFFFF7F,
	}

	for in, exp := range m {
		r := encodeLength(in)
		if r != exp {
			t.Errorf("encodeLength failed, input %d expected 0x%X, got 0x%X", in, exp, r)
		}
	}
}

func Test_decodeRemLen(t *testing.T) {
	m := map[uint32]uint{
		0x00000000: 0,
		0x0000007F: 127,
		0x00008001: 128,
		0x0000FF7F: 16383,
		0x00808001: 16384,
		0x00FFFF7F: 2097151,
		0x80808001: 2097152,
		0xFFFFFF7F: 268435455,
	}

	for in, exp := range m {
		r := decodeRemLen(in)
		if r != exp {
			t.Errorf("decodeRemLen failed, input 0x%x expected %d, got %d", in, exp, r)
		}
	}
}

func Test_NewConnectMessage_vheader(t *testing.T) {

	cm := newConnectMsg(false, false, QOS_ZERO, false, "/wills", []byte("mywill"), "mycid", "", "", 0)

	if cm.QoS() != QOS_ZERO {
		t.Errorf("NewConnectMessage wrong QoS")
	}
	if string(cm.vheader[2:8]) != "MQIsdp" {
		t.Errorf("NewConnectMessage wrong protocol version name, expected MQIsdp, got %s", string(cm.vheader[2:8]))
	}
	if cm.vheader[8] != 0x03 {
		t.Errorf("NewConnectMessage wrong protocol version number, expected 0x3, got 0x%X", cm.vheader[8])
	}
	if cm.vheader[9] != 0x00 {
		t.Errorf("NewConnectMessage bad connect byte, expected 0x0, got 0x%X", cm.vheader[9])
	}
	timeout := cm.timeout()
	if timeout != 0 {
		t.Errorf("NewConnectMessage bad timeout, expected 0, got %d", timeout)
	}

	cm = newConnectMsg(true, true, QOS_ONE, true, "/wills", []byte("mywill"), "mycid", "myself", "mypass", 15)
	if cm.QoS() != QOS_ZERO {
		t.Errorf("NewConnectMessage wrong Qos")
	}
	if string(cm.vheader[2:8]) != "MQIsdp" {
		t.Errorf("NewConnectMessage wrong protocol version name, expected MQIsdp, got %s", string(cm.vheader[2:8]))
	}
	if cm.vheader[8] != 0x03 {
		t.Errorf("NewConnectMessage wrong protocol version number, expected 0x3, got 0x%X", cm.vheader[8])
	}
	if cm.vheader[9] != 0xEE {
		t.Errorf("NewConnectMessage bad connect byte, expected 0xEE, got 0x%X", cm.vheader[9])
	}
	timeout = cm.timeout()
	if timeout != 15 {
		t.Errorf("NewConnectMessage bad timeout, expected 15, got %d", timeout)
	}

	cm = newConnectMsg(true, false, QOS_TWO, false, "/wills", []byte("mywill"), "mycid", "myuser", "", 27)
	if cm.QoS() != QOS_ZERO {
		t.Errorf("NewConnectMessage wrong Qos")
	}
	if string(cm.vheader[2:8]) != "MQIsdp" {
		t.Errorf("NewConnectMessage wrong protocol version name, expected MQIsdp, got %s", string(cm.vheader[2:8]))
	}
	if cm.vheader[8] != 0x03 {
		t.Errorf("NewConnectMessage wrong protocol version number, expected 0x3, got 0x%X", cm.vheader[8])
	}
	if cm.vheader[9] != 0x92 {
		t.Errorf("NewConnectMessage bad connect byte, expected 0x92, got 0x%X", cm.vheader[9])
	}
	timeout = cm.timeout()
	if timeout != 27 {
		t.Errorf("NewConnectMessage bad timeout, expected 15, got %d", timeout)
	}

	cm = newConnectMsg(true, false, QOS_TWO, false, "/wills", []byte("mywill"), "mycid", "myuser", "pass", 0xAB21)
	timeout = cm.timeout()
	if timeout != 0xAB21 {
		t.Errorf("NewConnectMessage bad timeout, expected 15, got %d", timeout)
	}
}

func Test_newConnectMsg_payload(t *testing.T) {
	cm := newConnectMsg(true, false, QOS_TWO, false, "/wills", []byte("mywill"), "mycid", "", "", 0xAB21)
	if cm.remLen() != 19 {
		t.Errorf("NewConnectionMessage bad remlen, expected 19, got %d", cm.remLen())
	}

	cm = newConnectMsg(true, true, QOS_TWO, false, "/wills", []byte("mywill"), "mycid", "", "", 0xAB21)
	if cm.remLen() != 35 {
		t.Errorf("NewConnectionMessage bad remlen, expected 35, got %d", cm.remLen())
	}
	if cm.payload[0] != 0x00 || cm.payload[1] != 0x05 { // len($clientid)
		t.Errorf("NewConnectionMessage bad will topic length, expected 6, got %d", cm.payload[1])
	}
	if string(cm.payload[2:7]) != "mycid" {
		t.Errorf("NewConnectionMessage bad client id, expected mycid, got %s", string(cm.payload[2:7]))
	}
	if cm.payload[7] != 0x00 || cm.payload[8] != 0x06 { // len($willtopic)
		t.Errorf("NewConnectionMessage bad willtopic length, expected 6, got %d", cm.payload[8])
	}
	if string(cm.payload[9:15]) != "/wills" {
		t.Errorf("NewConnectionMessage bad willtopic, expected /wills, got %s", string(cm.payload[9:15]))
	}
	if cm.payload[15] != 0x00 || cm.payload[16] != 0x06 { // len($willmsg)
		t.Errorf("NewConnectionMessage bad will message length, expected 6, got %d", cm.payload[16])
	}
	if string(cm.payload[17:23]) != "mywill" {
		t.Errorf("NewConnectionMessage bad will message, expected mywill, got %s", string(cm.payload[17:23]))
	}

	cm = newConnectMsg(true, true, QOS_TWO, false, "/w", []byte("will"), "cid", "User", "Password", 0x1234)

	if cm.remLen() != 43 {
		t.Errorf("NewConnectionMessage bad remlen, expected 35, got %d", cm.remLen())
	}
	if cm.payload[0] != 0x00 || cm.payload[1] != 0x03 { // len($clientid)
		t.Errorf("NewConnectionMessage bad will topic length, expected 6, got %d", cm.payload[1])
	}
	if string(cm.payload[2:5]) != "cid" {
		t.Errorf("NewConnectionMessage bad client id, expected mycid, got %s", string(cm.payload[2:5]))
	}
	if cm.payload[5] != 0x00 || cm.payload[6] != 0x02 { // len($willtopic)
		t.Errorf("NewConnectionMessage bad willtopic length, expected 6, got %d", cm.payload[6])
	}
	if string(cm.payload[7:9]) != "/w" {
		t.Errorf("NewConnectionMessage bad willtopic, expected /wills, got %s", string(cm.payload[7:9]))
	}
	if cm.payload[9] != 0x00 || cm.payload[10] != 0x04 { // len($willmsg)
		t.Errorf("NewConnectionMessage bad will message length, expected 6, got %d", cm.payload[10])
	}
	if string(cm.payload[11:15]) != "will" {
		t.Errorf("NewConnectionMessage bad will message, expected will, got %s", string(cm.payload[11:15]))
	}
	if cm.payload[15] != 0 || cm.payload[16] != 0x04 { // len($username)
		t.Errorf("NewConnectionMessage bad username length, expected 4, got %d", cm.payload[16])
	}
	if string(cm.payload[17:21]) != "User" {
		t.Errorf("NewConnectionMessage bad username, expected User, got %s", string(cm.payload[17:21]))
	}
	if cm.payload[21] != 0 || cm.payload[22] != 8 { // len($password)
		t.Errorf("NewConnectionMessage bad password length, expected 8, got %d", cm.payload[22])
	}
	if string(cm.payload[23:31]) != "Password" {
		t.Errorf("NewConnectionMessage bad password, expected Password, got %s", string(cm.payload[23:31]))
	}
}

func Test_Disconnect_Message(t *testing.T) {
	dm := newDisconnectMsg()
	if dm.remLen() != 0 {
		t.Fatalf("NewDisconnectMessage bad remlen, expected 0, got %d", dm.remLen())
	}
	if dm.msgType() != DISCONNECT {
		t.Fatalf("NewDisconnectMessage bad msg type, got %v", dm.msgType())
	}

	bs := dm.Bytes()
	exp := []byte{
		0xE0,
		0x00,
	}

	if len(bs) != 2 {
		t.Fatalf("len(Disconnect Msg Bytes) != 2, got: %d", bs)
	}
	if bs[0] != exp[0] {
		t.Fatalf("firt byte of Disconnect message is bad, got: 0x%X", bs[0])
	}

	if bs[1] != exp[1] {
		t.Fatalf("second byte of Disconnect message is bad, got: 0x%X", bs[1])
	}
}

func Test_newPublishMessage_pm0(t *testing.T) {
	pm := newPublishMsg(QOS_ZERO, "/a/b/c", []byte{0xBE, 0xEF, 0xED})
	exp := []byte{
		/* msg type */
		0x30,

		/* remlen */
		0x0B,

		/* topic, msg id in varheader */
		0x00, // length of topic
		0x06,
		0x2F, // /
		0x61, // a
		0x2F, // /
		0x62, // b
		0x2F, // /
		0x63, // c
		/* no message id, because qos 0 */

		/*payload */
		0xBE,
		0xEF,
		0xED,
	}
	bs := pm.Bytes()

	if len(bs) != len(exp) {
		t.Fatalf("new publish message has wrong number of bytes: %d, expected: %d", len(bs), len(exp))
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("new publish message has bad byte %x, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_newPublishMessage_pm1(t *testing.T) {
	pm := newPublishMsg(QOS_ONE, "/a/b/c", []byte{0xBE, 0xEF, 0xED})
	pm.setMsgId(9080)
	exp := []byte{
		/* msg type */
		0x32,

		/* remlen */
		0x0D,

		/* topic, msg id in varheader */
		0x00, // length of topic
		0x06,
		0x2F, // /
		0x61, // a
		0x2F, // /
		0x62, // b
		0x2F, // /
		0x63, // c

		/* message id  9080 base 10 */
		0x23,
		0x78,

		/* payload */
		0xBE,
		0xEF,
		0xED,
	}
	bs := pm.Bytes()

	if len(bs) != len(exp) {
		t.Fatalf("new publish message has wrong number of bytes: %d, expected: %d", len(bs), len(exp))
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("new publish message has bad byte #%d, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_newPubRelMessage(t *testing.T) {
	prm := newPubRelMsg()
	prm.setMsgId(46)
	exp := []byte{
		/* msg type */
		0x62,

		/* remlen */
		0x02,

		/* msg id 46 base 10 */
		0x00,
		0x2E,

		/* no payload */
	}

	bs := prm.Bytes()
	if len(bs) != len(exp) {
		t.Fatalf("new pubrel message has wrong number of bytes: %d, expected: %d", len(bs), len(exp))
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("new pubrel message has bad byte #%d, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_newSubscribeMessage(t *testing.T) {
	filter, _ := NewTopicFilter("/a", 0)
	sm := newSubscribeMsg(filter)
	sm.setMsgId(32)
	bs := sm.Bytes()
	exp := []byte{
		/* msg type */
		0x82,

		/* remlen */
		0x07,

		/* vheader (msg id) */
		0x00,
		0x20,

		/* payload */
		0x00, // topic length
		0x02, // topic length
		0x2F, // /
		0x61, // a
		0x00, // requested qos
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("new subscribe message has bad byte #%d, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_newSubscribeMessage_multi(t *testing.T) {
	filter1, _ := NewTopicFilter("/a", 1)
	filter2, _ := NewTopicFilter("/b", 0)
	filter3, _ := NewTopicFilter("/x/y/z", 2)
	sm := newSubscribeMsg(filter1, filter2, filter3)
	sm.setMsgId(3432)
	if sm == nil {
		t.Fatalf("newSubscribe message is nil")
	}

	exp := []byte{
		/* msg type */
		0x82,

		/* remlen */
		0x15,

		/* vheader (msg id) */
		0x0D,
		0x68,

		/* payload */
		0x00, // topic length
		0x02, // topic length
		0x2F, // /
		0x61, // a
		0x01, // requested qos

		0x00, // topic length
		0x02, // topic length
		0x2F, // /
		0x62, // b
		0x00, // requested qos

		0x00, // topic length
		0x06, // topic length
		0x2F, // /
		0x78, // x
		0x2F, // /
		0x79, // y
		0x2F, // /
		0x7A, // z
		0x02, // requested qos
	}
	bs := sm.Bytes()

	if len(bs) != len(exp) {
		t.Fatalf("newSubscribe message has wrong number of bytes, expected %d got %d", len(exp), len(bs))
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("new subscribe message has bad byte #%d, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_newUnsubscribeMessage_unsb1(t *testing.T) {
	um := newUnsubscribeMsg("z")
	um.setMsgId(4213)

	exp := []byte{
		0xA2, // msg type
		0x05, // remlen
		0x10, // msg id (msb)
		0x75, // msg id (lsb)
		0x00, // topic length
		0x01, // topic length
		0x7A,
	}

	bs := um.Bytes()
	if len(bs) != len(exp) {
		t.Fatalf("newUnsubscribe message has wrong number of bytes, expected %d got %d", len(exp), len(bs))
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("newUnsubscribe message has bad byte #%d, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_newUnsubscribeMessage_unsb3(t *testing.T) {
	um := newUnsubscribeMsg("/a", "/a/b", "/a/b/c")
	um.setMsgId(4213)

	exp := []byte{
		0xA2, // msg type
		0x14, // remlen

		0x10, // msg id (msb)
		0x75, // msg id (lsb)

		0x00, // topic length (1)
		0x02, // topic length (1)
		0x2F, // /
		0x61, // a

		0x00, // topic length (2)
		0x04, // topic length (2)
		0x2F, // /
		0x61, // a
		0x2F, // /
		0x62, // b

		0x00, // topic length (3)
		0x06, // topic length (3)
		0x2F, // /
		0x61, // a
		0x2F, // /
		0x62, // b
		0x2F, // /
		0x63, // c
	}

	bs := um.Bytes()
	if len(bs) != len(exp) {
		t.Fatalf("newUnsubscribe message has wrong number of bytes, expected %d got %d", len(exp), len(bs))
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("newUnsubscribe message has bad byte #%d, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_NewMessage(t *testing.T) {
	m := NewMessage([]byte("test message"))

	exp := []byte{
		/* msg type */
		0x32,
		/* remlen */
		0x00, // remlen - not used just yet
		0x74, // t
		0x65, // e
		0x73, // s
		0x74, // t
		0x20, // ' '
		0x6D, // m
		0x65, // e
		0x73, // s
		0x73, // s
		0x61, // a
		0x67, // g
		0x65, // e
	}
	bs := m.Bytes()

	if len(bs) != len(exp) {
		t.Fatalf("new publish message has wrong number of bytes: %d, expected: %d", len(bs), len(exp))
	}

	for i := range exp {
		if bs[i] != exp[i] {
			t.Fatalf("new publish message has bad byte #%d, expected 0x%X, got 0x%X", i, exp[i], bs[i])
		}
	}
}

func Test_Topic_from_publish(t *testing.T) {
	pm := newPublishMsg(QOS_ZERO, "/a/b/c", []byte{0xBE, 0xEF, 0xED})

	/*
		// msg type
		0x30

		// remlen
		0x0B

		// topic, msg id in varheader
		0x00 // length of topic
		0x06
		0x2F // /
		0x61 // a
		0x2F // /
		0x62 // b
		0x2F // /
		0x63 // c
		// no message id, because qos 0

		// payload
		0xBE
		0xEF
		0xED
	*/

	topic := pm.Topic()

	if topic != "/a/b/c" {
		t.Fatalf("getTopic is incorrect")
	}
}

func Test_set_cleanSession(t *testing.T) {
	ops := NewClientOptions().SetCleanSession(false)
	c := NewClient(ops)
	if c.options.cleanSession {
		t.Fatalf("cleanSession was true but was set false")
	}
}
