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

import "bytes"
import "testing"

func Test_decodeQos_2(t *testing.T) {
	b := byte(0xFD) // "4" is qos 2
	qos := decodeQos(b)
	if qos != QOS_TWO {
		t.Fatalf("decod_qos expected 2, got %d", qos)
	}
}

func Test_decodeQos_1(t *testing.T) {
	b := byte(0xFB) // "2" is qos 1
	qos := decodeQos(b)
	if qos != QOS_ONE {
		t.Fatalf("decod_qos expected 1, got %d", qos)
	}
}

func Test_decodeQos_0(t *testing.T) {
	b := byte(0xF9) // "0" is qos 0
	qos := decodeQos(b)
	if qos != QOS_ZERO {
		t.Fatalf("decodeQos expected 0, got %d", qos)
	}
}

func Test_decodeMsgType_connect(t *testing.T) {
	b := byte(0x1F)
	mtype := decodeMsgType(b)
	if mtype != CONNECT {
		t.Fatalf("decodeMsgType expected %d, got %d", CONNECT, mtype)
	}
}

func helper_test_decodeRemlen(t *testing.T, bs []byte, exp_n int, exp_v uint32) {
	n, v := decodeRemlen(bs)
	if n != exp_n {
		t.Fatalf("decodeRemlen n expected: %d, got %d", exp_n, n)
	}
	if v != exp_v {
		t.Fatalf("decodeRemlen v expected: %d, got %d", exp_v, v)
	}
}

func Test_decodeRemlen_1_low(t *testing.T) {
	bs := []byte{
		0xAB,
		0x00, // 1 remlen byte (value 0)
	}
	helper_test_decodeRemlen(t, bs, 1, 0)
}

func Test_decodeRemlen_1_high(t *testing.T) {
	bs := []byte{
		0xAB,
		0x7F, // 1 remlen byte (value 127)
	}
	helper_test_decodeRemlen(t, bs, 1, 127)
}

func Test_decodeRemlen_2_low(t *testing.T) {
	bs := []byte{
		0xAB,
		0x80, // 2 remlen byte (value 128)
		0x01,
	}
	helper_test_decodeRemlen(t, bs, 2, 128)
}

func Test_decodeRemlen_2_high(t *testing.T) {
	bs := []byte{
		0xAB,
		0xFF, // 2 remlen byte (value 16383)
		0x7F,
	}
	helper_test_decodeRemlen(t, bs, 2, 16383)
}

func Test_decodeRemlen_3_low(t *testing.T) {
	bs := []byte{
		0xAB,
		0x80, // 3 remlen byte (value 16384)
		0x80,
		0x01,
	}
	helper_test_decodeRemlen(t, bs, 3, 16384)
}

func Test_decodeRemlen_3_high(t *testing.T) {
	bs := []byte{
		0xAB,
		0xFF, // 3 remlen byte (value 2097151)
		0xFF,
		0x7F,
	}
	helper_test_decodeRemlen(t, bs, 3, 2097151)
}

func Test_decodeRemlen_4_low(t *testing.T) {
	bs := []byte{
		0xAB,
		0x80, // 4 remlen byte (value 2097152)
		0x80,
		0x80,
		0x01,
	}
	helper_test_decodeRemlen(t, bs, 4, 2097152)
}

func Test_decodeRemlen_4_high(t *testing.T) {
	bs := []byte{
		0xAB,
		0xFF, // 4 remlen byte (value 268435455)
		0xFF,
		0xFF,
		0x7F,
	}
	helper_test_decodeRemlen(t, bs, 4, 268435455)
}

func Test_decodeTopic(t *testing.T) {
	bs := []byte{
		0x00,
		0x03,
		0x41,
		0x42,
		0x43,
	}
	n, s := decodeTopic(bs)
	if n != 3 {
		t.Fatalf("decod_topic n expected %d, got %d", 3, n)
	}
	if s != "ABC" {
		t.Fatalf("decodeTopic t expected %s, got %s", "ABC", s)
	}
}

func Test_decode_connack(t *testing.T) {
	connack := []byte{
		/* Fixed Header */
		0x20, // msgtype (CONNACK)
		0x02, // remlen
		/* Variable Header */
		0x00, // reserved
		0x00, // return code (CONN_ACCEPTED)
	}
	m := decode(connack)
	if m.msgType() != CONNACK {
		t.Errorf("decode bad message type expected CONNACK, got %v", m.msgType())
	}
	if m.remLen() != 2 {
		t.Errorf("decode bad remaining length expected 2, got %d", m.remLen())
	}
	if m.connRC() != CONN_ACCEPTED {
		t.Errorf("decode badd conn rc expected CONN_ACCEPTED, got %v", m.connRC())
	}

	connack = []byte{
		/* Fixed Header */
		0x20, // msgtype (CONNACK)
		0x02, // remlen
		/* Variable Header */
		0x00, // reserved
		0x01, // return code (CONN_REF_BAD_PROTO_VER)
	}
	m = decode(connack)
	if m.msgType() != CONNACK {
		t.Errorf("decode bad message type expected CONNACK, got %v", m.msgType())
	}
	if m.remLen() != 2 {
		t.Errorf("decode bad remaining length expected 2, got %d", m.remLen())
	}
	if m.connRC() != CONN_REF_BAD_PROTO_VER {
		t.Errorf("decode badd conn rc expected CONN_REF_BAD_PROTO_VER, got %v", m.connRC())
	}

	connack = []byte{
		/* Fixed Header */
		0x20, // msgtype (CONNACK)
		0x02, // remlen
		/* Variable Header */
		0x00, // reserved
		0x02, // return code (CONN_REF_ID_REJ
	}
	m = decode(connack)
	if m.msgType() != CONNACK {
		t.Errorf("decode bad message type expected CONNACK, got %v", m.msgType())
	}
	if m.remLen() != 2 {
		t.Errorf("decode bad remaining length expected 2, got %d", m.remLen())
	}
	if m.connRC() != CONN_REF_ID_REJ {
		t.Errorf("decode badd conn rc expected CONN_REF_ID_REJ, got %v", m.connRC())
	}

	connack = []byte{
		/* Fixed Header */
		0x20, // msgtype (CONNACK)
		0x02, // remlen
		/* Variable Header */
		0x00, // reserved
		0x03, // return code (CONN_REF_SERV_UNAVAIL)
	}
	m = decode(connack)
	if m.msgType() != CONNACK {
		t.Errorf("decode bad message type expected CONNACK, got %v", m.msgType())
	}
	if m.remLen() != 2 {
		t.Errorf("decode bad remaining length expected 2, got %d", m.remLen())
	}
	if m.connRC() != CONN_REF_SERV_UNAVAIL {
		t.Errorf("decode badd conn rc expected CONN_REF_SERV_UNAVAIL, got %v", m.connRC())
	}

	connack = []byte{
		/* Fixed Header */
		0x20, // msgtype (CONNACK)
		0x02, // remlen
		/* Variable Header */
		0x00, // reserved
		0x04, // return code (CONN_REF_BAD_USER_PASS
	}
	m = decode(connack)
	if m.msgType() != CONNACK {
		t.Errorf("decode bad message type expected CONNACK, got %v", m.msgType())
	}
	if m.remLen() != 2 {
		t.Errorf("decode bad remaining length expected 2, got %d", m.remLen())
	}
	if m.connRC() != CONN_REF_BAD_USER_PASS {
		t.Errorf("decode badd conn rc expected CONN_REF_BAD_USER_PASS got %v", m.connRC())
	}

	connack = []byte{
		/* Fixed Header */
		0x20, // msgtype (CONNACK)
		0x02, // remlen
		/* Variable Header */
		0x00, // reserved
		0x05, // return code (CONN_REF_NOT_AUTH)
	}
	m = decode(connack)
	if m.msgType() != CONNACK {
		t.Errorf("decode bad message type expected CONNACK, got %v", m.msgType())
	}
	if m.remLen() != 2 {
		t.Errorf("decode bad remaining length expected 2, got %d", m.remLen())
	}
	if m.connRC() != CONN_REF_NOT_AUTH {
		t.Errorf("decode badd conn rc expected CONN_REF_NOT_AUTH, got %v", m.connRC())
	}
}

func Benchmark_decode_connack(b *testing.B) {
	connack := []byte{
		/* Fixed Header */
		0x20, // msgtype (CONNACK)
		0x02, // remlen
		/* Variable Header */
		0x00, // reserved
		0x04, // return code (CONN_REF_BAD_USER_PASS
	}
	for i := 0; i < b.N; i++ {
		decode(connack)
	}
}

func Test_Bytes_connect(t *testing.T) {
	m := newConnectMsg(false, false, QOS_ZERO, false, "", []byte(""), "mycid", "", "", 10)
	bs := m.Bytes()
	if len(bs) != 21 {
		t.Errorf("len(m.Bytes()) is wrong, expected 21, got %d", len(bs))
	}
	exp := []byte{
		0x10, // msgtype, dup, qos
		0x13, // remlen

		0x00, // strlen
		0x06, // strlen
		'M',  // proto name
		'Q',
		'I',
		's',
		'd',
		'p',

		0x03, // proto version

		0x00, // connect flags

		0x00, // keep alive
		0x0A, // keep alive

		0x00, // strlen
		0x05,
		'm',
		'y',
		'c',
		'i',
		'd',
	}

	for i := range exp {
		if exp[i] != bs[i] {
			t.Errorf("exp[%d] != bs[%d], 0x%X != 0x%X", i, i, exp[i], bs[i])
		}
	}

	m = newConnectMsg(true, true, QOS_TWO, true, "/all/wills/go/here", []byte("good bye everybody :("), "TheClientID", "TheUserName", "ThePassWord", 0xABCD)
	bs = m.Bytes()
	if len(bs) != 96 {
		t.Errorf("len(m.Bytes()) is wrong, expected 96, got %d", len(bs))
	}
	exp = []byte{
		0x10, // msgtype, dup, qos
		0x5E, // remlen

		0x00, // strlen
		0x06, // strlen
		'M',  // proto name
		'Q',
		'I',
		's',
		'd',
		'p',

		0x03, // proto version

		0xF6, // connect flags

		0xAB, // keep alive
		0xCD, // keep alive

		0x00, // client id
		0x0B,
		'T',
		'h',
		'e',
		'C',
		'l',
		'i',
		'e',
		'n',
		't',
		'I',
		'D',

		0x00, // will topic
		0x12,
		'/',
		'a',
		'l',
		'l',
		'/',
		'w',
		'i',
		'l',
		'l',
		's',
		'/',
		'g',
		'o',
		'/',
		'h',
		'e',
		'r',
		'e',

		0x00, // will msg
		0x15,
		'g',
		'o',
		'o',
		'd',
		' ',
		'b',
		'y',
		'e',
		' ',
		'e',
		'v',
		'e',
		'r',
		'y',
		'b',
		'o',
		'd',
		'y',
		' ',
		':',
		'(',

		0x00, // usrname
		0x0B,
		'T',
		'h',
		'e',
		'U',
		's',
		'e',
		'r',
		'N',
		'a',
		'm',
		'e',

		0x00, // password
		0x0B,
		'T',
		'h',
		'e',
		'P',
		'a',
		's',
		's',
		'W',
		'o',
		'r',
		'd',
	}

	for i := range exp {
		if exp[i] != bs[i] {
			t.Errorf("exp[%d] != bs[%d], 0x%X != 0x%X", i, i, exp[i], bs[i])
		}
	}
}

func Test_decode_puback(t *testing.T) {
	bs := []byte{
		0x40, // type
		0x02, // remlen
		0x00, // msg id (msb)
		0x09, // msg id (lsb)
	}

	m := decode(bs)

	if m.msgType() != PUBACK {
		t.Errorf("decode puback wrong message type: %v", m.msgType())
	}

	if m.remLen() != 2 {
		t.Errorf("decode puback wrong remlen: %v", m.remLen())
	}

	if m.MsgId() != 9 {
		t.Errorf("decode puback wrong message id: %v", m.MsgId())
	}
}

func Test_decode_suback(t *testing.T) {
	bs := []byte{
		0x90, // type
		0x00, // remlen
		0xAB, // msg id (msb)
		0x12, // msg id (lsb)
		0x01, // qos (1)
		0x02, // qos (2)
		0x00, // qos (3)
		0x00, // qos (4)
		0x01, // qos (5)
	}
	m := decode(bs)
	if m.msgType() != SUBACK {
		t.Errorf("decode suback wrong message type %v", m.msgType())
	}
	if m.remLen() != 0 {
		t.Errorf("decode suback wrong remlen: %d", m.remLen())
	}
	if m.MsgId() != 43794 {
		t.Fatalf("decode suback wrong message id %d", m.MsgId())
	}
	if !bytes.Equal(m.Payload(), []byte{1, 2, 0, 0, 1}) {
		t.Fatalf("decode wrong payload")
	}
}

func Test_decode_unsuback(t *testing.T) {
	bs := []byte{
		0xB0, // type
		0x00, // remlen
		0xAB, // msg id (msb)
		0x12, // msg id (lsb)
	}
	m := decode(bs)
	if m.msgType() != UNSUBACK {
		t.Errorf("decode unsuback wrong message type %v", m.msgType())
	}
	if m.remLen() != 0 {
		t.Errorf("decode unsuback wrong remlen: %d", m.remLen())
	}
	if m.MsgId() != 43794 {
		t.Fatalf("decode unsuback wrong message id %d", m.MsgId())
	}
}
