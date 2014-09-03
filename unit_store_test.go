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
	"fmt"
	. "github.com/alsm/hrotti/packets"
	"io/ioutil"
	"os"
	"testing"
)

func Test_fullpath(t *testing.T) {
	p := fullpath("/tmp/store", "o.44324")
	e := "/tmp/store/o.44324.msg"
	if p != e {
		t.Fatalf("full path expected %s, got %s", e, p)
	}
}

func Test_exists(t *testing.T) {
	b := exists("/")
	if !b {
		t.Errorf("/proc/cpuinfo was not found")
	}
}

func Test_exists_no(t *testing.T) {
	b := exists("/this/path/is/not/real/i/hope")
	if b {
		t.Errorf("you have some strange files")
	}
}

func isemptydir(dir string) bool {
	chkcond(exists(dir))
	files, err := ioutil.ReadDir(dir)
	chkerr(err)
	return len(files) == 0
}

func Test_key2mid(t *testing.T) {
	key := "i.123"
	exp := uint16(123)
	res := key2mid(key)
	if exp != res {
		t.Fatalf("key2mid failed")
	}
}

func Test_ibound_mid2key(t *testing.T) {
	id := uint16(9876)
	exp := "i.9876"
	res := ibound_mid2key(id)
	if exp != res {
		t.Fatalf("ibound_mid2key failed")
	}
}

func Test_obound_mid2key(t *testing.T) {
	id := uint16(7654)
	exp := "o.7654"
	res := obound_mid2key(id)
	if exp != res {
		t.Fatalf("obound_mid2key failed")
	}
}

/************************
 **** persist_obound ****
 ************************/

func Test_persist_obound_connect(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(CONNECT).(*ConnectPacket)
	m.Qos = 0
	m.Username = "user"
	m.Password = []byte("pass")
	m.ClientIdentifier = "cid"
	//m := newConnectMsg(false, false, QOS_ZERO, false, "", nil, "cid", "user", "pass", 10)
	persist_obound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_publish_0(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBLISH).(*PublishPacket)
	m.Qos = 0
	m.TopicName = "/popub0"
	m.Payload = []byte{0xBB, 0x00}
	m.MessageID = 40
	persist_obound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_publish_1(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBLISH).(*PublishPacket)
	m.Qos = 1
	m.TopicName = "/popub1"
	m.Payload = []byte{0xBB, 0x00}
	m.MessageID = 41
	persist_obound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 41 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_publish_2(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBLISH).(*PublishPacket)
	m.Qos = 2
	m.TopicName = "/popub2"
	m.Payload = []byte{0xBB, 0x00}
	m.MessageID = 42
	persist_obound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 42 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_puback(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBACK).(*PubackPacket)
	persist_obound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 1 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_pubrec(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBREC).(*PubrecPacket)
	persist_obound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_pubrel(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBREL).(*PubrelPacket)
	m.MessageID = 43

	persist_obound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 43 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_pubcomp(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBCOMP).(*PubcompPacket)
	persist_obound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 1 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_subscribe(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(SUBSCRIBE).(*SubscribePacket)
	m.Topics = []string{"/posub"}
	m.Qoss = []byte{1}
	m.MessageID = 44
	persist_obound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 44 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_unsubscribe(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(UNSUBSCRIBE).(*UnsubscribePacket)
	m.Topics = []string{"/posub"}
	m.MessageID = 45
	persist_obound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 45 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_pingreq(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PINGREQ)
	persist_obound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

func Test_persist_obound_disconnect(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(DISCONNECT)
	persist_obound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_obound put message it should not have")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_obound get message it should not have")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_obound del message it should not have")
	}
}

/************************
 **** persist_ibound ****
 ************************/

func Test_persist_ibound_connack(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(CONNACK)
	persist_ibound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_publish_0(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBLISH).(*PublishPacket)
	m.Qos = 0
	m.TopicName = "/pipub0"
	m.Payload = []byte{0xCC, 0x01}
	m.MessageID = 50
	persist_ibound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_publish_1(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBLISH).(*PublishPacket)
	m.Qos = 1
	m.TopicName = "/pipub1"
	m.Payload = []byte{0xCC, 0x02}
	m.MessageID = 51
	persist_ibound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 51 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_publish_2(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PUBLISH).(*PublishPacket)
	m.Qos = 2
	m.TopicName = "/pipub2"
	m.Payload = []byte{0xCC, 0x03}
	m.MessageID = 52
	persist_ibound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 52 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_puback(t *testing.T) {
	ts := &TestStore{}
	pub := NewControlPacket(PUBLISH).(*PublishPacket)
	pub.Qos = 1
	pub.TopicName = "/pub1"
	pub.Payload = []byte{0xCC, 0x04}
	pub.MessageID = 53
	publish_key := ibound_mid2key(pub.MessageID)
	ts.Put(publish_key, pub)

	m := NewControlPacket(PUBACK).(*PubackPacket)
	m.MessageID = 53

	persist_ibound(ts, m) // "deletes" PUBLISH from store

	if len(ts.mput) != 1 { // not actually deleted in TestStore
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 1 || ts.mdel[0] != 53 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_pubrec(t *testing.T) {
	ts := &TestStore{}
	pub := NewControlPacket(PUBLISH).(*PublishPacket)
	pub.Qos = 2
	pub.TopicName = "/pub2"
	pub.Payload = []byte{0xCC, 0x05}
	pub.MessageID = 54
	publish_key := ibound_mid2key(pub.MessageID)
	ts.Put(publish_key, pub)

	m := NewControlPacket(PUBREC).(*PubrecPacket)
	m.MessageID = 54

	persist_ibound(ts, m)

	if len(ts.mput) != 1 || ts.mput[0] != 54 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_pubrel(t *testing.T) {
	ts := &TestStore{}
	pub := NewControlPacket(PUBLISH).(*PublishPacket)
	pub.Qos = 2
	pub.TopicName = "/pub2"
	pub.Payload = []byte{0xCC, 0x06}
	pub.MessageID = 55
	publish_key := ibound_mid2key(pub.MessageID)
	ts.Put(publish_key, pub)

	m := NewControlPacket(PUBREL).(*PubrelPacket)
	m.MessageID = 55

	persist_ibound(ts, m) // will overwrite publish

	if len(ts.mput) != 2 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_pubcomp(t *testing.T) {
	ts := &TestStore{}

	m := NewControlPacket(PUBCOMP).(*PubcompPacket)
	m.MessageID = 56

	persist_ibound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 1 || ts.mdel[0] != 56 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_suback(t *testing.T) {
	ts := &TestStore{}

	m := NewControlPacket(SUBACK).(*SubackPacket)
	m.MessageID = 57

	persist_ibound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 1 || ts.mdel[0] != 57 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_unsuback(t *testing.T) {
	ts := &TestStore{}

	m := NewControlPacket(UNSUBACK).(*UnsubackPacket)
	m.MessageID = 58

	persist_ibound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 1 || ts.mdel[0] != 58 {
		t.Fatalf("persist_ibound in bad state")
	}
}

func Test_persist_ibound_pingresp(t *testing.T) {
	ts := &TestStore{}
	m := NewControlPacket(PINGRESP)

	persist_ibound(ts, m)

	if len(ts.mput) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mget) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}

	if len(ts.mdel) != 0 {
		t.Fatalf("persist_ibound in bad state")
	}
}

/***********
 * restore *
 ***********/

func ensure_restore_dir() {
	if exists("/tmp/restore") {
		rerr := os.RemoveAll("/tmp/restore")
		chkerr(rerr)
	}
	os.Mkdir("/tmp/restore", 0766)
}

func write_to_restore(fname, content string) {
	f, cerr := os.Create("/tmp/restore/" + fname)
	chkerr(cerr)
	chkcond(f != nil)
	w := bufio.NewWriter(f)
	w.Write([]byte(content))
	w.Flush()
	f.Close()
}

func verify_from_restore(fname, content string, t *testing.T) {
	msg, oerr := os.Open("/tmp/restore/" + fname)
	chkerr(oerr)
	all, rerr := ioutil.ReadAll(msg)
	chkerr(rerr)
	msg.Close()
	s := string(all)
	if s != content {
		t.Fatalf("verify content expected `%s` but got `%s`")
	}
}

func Test_restore_1(t *testing.T) {
	ensure_restore_dir()

	write_to_restore("i.1.bkp", "this is critical 1")

	restore("/tmp/restore")

	chkcond(!exists("/tmp/restore/i.1.bkp"))
	chkcond(exists("/tmp/restore/i.1.msg"))

	verify_from_restore("i.1.msg", "this is critical 1", t)
}

func Test_restore_2(t *testing.T) {
	ensure_restore_dir()

	write_to_restore("o.2.msg", "this is critical 2")

	restore("/tmp/restore")

	chkcond(!exists("/tmp/restore/o.2.bkp"))
	chkcond(exists("/tmp/restore/o.2.msg"))

	verify_from_restore("o.2.msg", "this is critical 2", t)
}

func Test_restore_3(t *testing.T) {
	ensure_restore_dir()

	N := 20
	// evens are .msg
	// odds are .bkp
	for i := 0; i < N; i++ {
		content := fmt.Sprintf("foo %d bar", i)
		if i%2 == 0 {
			mname := fmt.Sprintf("i.%d.msg", i)
			write_to_restore(mname, content)
		} else {
			mname := fmt.Sprintf("i.%d.bkp", i)
			write_to_restore(mname, content)
		}
	}

	restore("/tmp/restore")

	for i := 0; i < N; i++ {
		mname := fmt.Sprintf("i.%d.msg", i)
		bname := fmt.Sprintf("i.%d.bkp", i)
		content := fmt.Sprintf("foo %d bar", i)
		chkcond(!exists("/tmp/restore/" + bname))
		chkcond(exists("/tmp/restore/" + mname))

		verify_from_restore(mname, content, t)
	}
}
