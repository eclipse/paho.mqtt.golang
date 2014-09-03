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
	"bytes"
	"fmt"
	. "github.com/alsm/hrotti/packets"
	"testing"
)

/*******************************
 **** Some helper functions ****
 *******************************/

func b2s(bs []byte) string {
	s := ""
	for _, b := range bs {
		s += fmt.Sprintf("%x ", b)
	}
	return s
}

/**********************************************
 **** A mock store implementation for test ****
 **********************************************/

type TestStore struct {
	mput []uint16
	mget []uint16
	mdel []uint16
}

func (ts *TestStore) Open() {
}

func (ts *TestStore) Close() {
}

func (ts *TestStore) Put(key string, m ControlPacket) {
	ts.mput = append(ts.mput, m.Details().MessageID)
}

func (ts *TestStore) Get(key string) ControlPacket {
	mid := key2mid(key)
	ts.mget = append(ts.mget, mid)
	return nil
}

func (ts *TestStore) All() []string {
	return nil
}

func (ts *TestStore) Del(key string) {
	mid := key2mid(key)
	ts.mdel = append(ts.mdel, mid)
}

func (ts *TestStore) Reset() {
}

/*******************
 **** FileStore ****
 *******************/

func Test_NewFileStore(t *testing.T) {
	storedir := "/tmp/TestStore/_new"
	f := NewFileStore(storedir)
	if f.opened {
		t.Fatalf("filestore was opened without opening it")
	}
	if f.directory != storedir {
		t.Fatalf("filestore directory is wrong")
	}
	// storedir might exist or might not, just like with a real client
	// the point is, we don't care, we just want it to exist after it is
	// opened
}

func Test_FileStore_Open(t *testing.T) {
	storedir := "/tmp/TestStore/_open"

	f := NewFileStore(storedir)
	f.Open()
	if !f.opened {
		t.Fatalf("filestore was not set open")
	}
	if f.directory != storedir {
		t.Fatalf("filestore directory is wrong")
	}
	if !exists(storedir) {
		t.Fatalf("filestore directory does not exst after opening it")
	}
}

func Test_FileStore_Close(t *testing.T) {
	storedir := "/tmp/TestStore/_unopen"
	f := NewFileStore(storedir)
	f.Open()
	if !f.opened {
		t.Fatalf("filestore was not set open")
	}
	if f.directory != storedir {
		t.Fatalf("filestore directory is wrong")
	}
	if !exists(storedir) {
		t.Fatalf("filestore directory does not exst after opening it")
	}

	f.Close()
	if f.opened {
		t.Fatalf("filestore was still open after unopen")
	}
	if !exists(storedir) {
		t.Fatalf("filestore was deleted after unopen")
	}
}

func Test_FileStore_write(t *testing.T) {
	storedir := "/tmp/TestStore/_write"
	f := NewFileStore(storedir)
	f.Open()

	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 1
	pm.TopicName = "a/b/c"
	pm.Payload = []byte{0xBE, 0xEF, 0xED}
	pm.MessageID = 91

	key := ibound_mid2key(pm.MessageID)
	f.Put(key, pm)

	if !exists(storedir + "/i.91.msg") {
		t.Fatalf("message not in store")
	}

}

func Test_FileStore_Get(t *testing.T) {
	storedir := "/tmp/TestStore/_get"
	f := NewFileStore(storedir)
	f.Open()
	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 1
	pm.TopicName = "/a/b/c"
	pm.Payload = []byte{0xBE, 0xEF, 0xED}
	pm.MessageID = 120

	key := obound_mid2key(pm.MessageID)
	f.Put(key, pm)

	if !exists(storedir + "/o.120.msg") {
		t.Fatalf("message not in store")
	}

	exp := []byte{
		/* msg type */
		0x32, // qos 1

		/* remlen */
		0x0d,

		/* topic, msg id in varheader */
		0x00, // length of topic
		0x06,
		0x2F, // /
		0x61, // a
		0x2F, // /
		0x62, // b
		0x2F, // /
		0x63, // c

		/* msg id (is always 2 bytes) */
		0x00,
		0x78,

		/*payload */
		0xBE,
		0xEF,
		0xED,
	}

	m := f.Get(key)

	if m == nil {
		t.Fatalf("message not retreived from store")
	}

	var msg bytes.Buffer
	m.Write(&msg)
	if !bytes.Equal(exp, msg.Bytes()) {
		t.Fatalf("message from store not same as what went in", msg.Bytes())
	}
}

func Test_FileStore_All(t *testing.T) {
	storedir := "/tmp/TestStore/_all"
	f := NewFileStore(storedir)
	f.Open()
	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 2
	pm.TopicName = "/t/r/v"
	pm.Payload = []byte{0x01, 0x02}
	pm.MessageID = 121

	key := obound_mid2key(pm.MessageID)
	f.Put(key, pm)

	keys := f.All()
	if len(keys) != 1 {
		t.Fatalf("FileStore.All does not have the messages")
	}

	if keys[0] != "o.121" {
		t.Fatalf("FileStore.All has wrong key")
	}
}

func Test_FileStore_Del(t *testing.T) {
	storedir := "/tmp/TestStore/_del"
	f := NewFileStore(storedir)
	f.Open()

	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 1
	pm.TopicName = "a/b/c"
	pm.Payload = []byte{0xBE, 0xEF, 0xED}
	pm.MessageID = 17

	key := ibound_mid2key(pm.MessageID)
	f.Put(key, pm)

	if !exists(storedir + "/i.17.msg") {
		t.Fatalf("message not in store")
	}

	f.Del(key)

	if exists(storedir + "/i.17.msg") {
		t.Fatalf("message still exists after deletion")
	}
}

func Test_FileStore_Reset(t *testing.T) {
	storedir := "/tmp/TestStore/_reset"
	f := NewFileStore(storedir)
	f.Open()

	pm1 := NewControlPacket(PUBLISH).(*PublishPacket)
	pm1.Qos = 1
	pm1.TopicName = "/q/w/e"
	pm1.Payload = []byte{0xBB}
	pm1.MessageID = 71
	key1 := ibound_mid2key(pm1.MessageID)
	f.Put(key1, pm1)

	pm2 := NewControlPacket(PUBLISH).(*PublishPacket)
	pm2.Qos = 1
	pm2.TopicName = "/q/w/e"
	pm2.Payload = []byte{0xBB}
	pm2.MessageID = 72
	key2 := ibound_mid2key(pm2.MessageID)
	f.Put(key2, pm2)

	pm3 := NewControlPacket(PUBLISH).(*PublishPacket)
	pm3.Qos = 1
	pm3.TopicName = "/q/w/e"
	pm3.Payload = []byte{0xBB}
	pm3.MessageID = 73
	key3 := ibound_mid2key(pm3.MessageID)
	f.Put(key3, pm3)

	pm4 := NewControlPacket(PUBLISH).(*PublishPacket)
	pm4.Qos = 1
	pm4.TopicName = "/q/w/e"
	pm4.Payload = []byte{0xBB}
	pm4.MessageID = 74
	key4 := ibound_mid2key(pm4.MessageID)
	f.Put(key4, pm4)

	pm5 := NewControlPacket(PUBLISH).(*PublishPacket)
	pm5.Qos = 1
	pm5.TopicName = "/q/w/e"
	pm5.Payload = []byte{0xBB}
	pm5.MessageID = 75
	key5 := ibound_mid2key(pm5.MessageID)
	f.Put(key5, pm5)

	if !exists(storedir + "/i.71.msg") {
		t.Fatalf("message not in store")
	}

	if !exists(storedir + "/i.72.msg") {
		t.Fatalf("message not in store")
	}

	if !exists(storedir + "/i.73.msg") {
		t.Fatalf("message not in store")
	}

	if !exists(storedir + "/i.74.msg") {
		t.Fatalf("message not in store")
	}

	if !exists(storedir + "/i.75.msg") {
		t.Fatalf("message not in store")
	}

	f.Reset()

	if exists(storedir + "/i.71.msg") {
		t.Fatalf("message still exists after reset")
	}

	if exists(storedir + "/i.72.msg") {
		t.Fatalf("message still exists after reset")
	}

	if exists(storedir + "/i.73.msg") {
		t.Fatalf("message still exists after reset")
	}

	if exists(storedir + "/i.74.msg") {
		t.Fatalf("message still exists after reset")
	}

	if exists(storedir + "/i.75.msg") {
		t.Fatalf("message still exists after reset")
	}
}

/*******************
 *** MemoryStore ***
 *******************/

func Test_NewMemoryStore(t *testing.T) {
	m := NewMemoryStore()
	if m == nil {
		t.Fatalf("MemoryStore could not be created")
	}
}

func Test_MemoryStore_Open(t *testing.T) {
	m := NewMemoryStore()
	m.Open()
	if !m.opened {
		t.Fatalf("MemoryStore was not set open")
	}
}

func Test_MemoryStore_Close(t *testing.T) {
	m := NewMemoryStore()
	m.Open()
	if !m.opened {
		t.Fatalf("MemoryStore was not set open")
	}

	m.Close()
	if m.opened {
		t.Fatalf("MemoryStore was still open after unopen")
	}
}

func Test_MemoryStore_Reset(t *testing.T) {
	m := NewMemoryStore()
	m.Open()

	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 2
	pm.TopicName = "/f/r/s"
	pm.Payload = []byte{0xAB}
	pm.MessageID = 81

	key := obound_mid2key(pm.MessageID)
	m.Put(key, pm)

	if len(m.messages) != 1 {
		t.Fatalf("message not in memstore")
	}

	m.Reset()

	if len(m.messages) != 0 {
		t.Fatalf("reset did not clear memstore")
	}
}

func Test_MemoryStore_write(t *testing.T) {
	m := NewMemoryStore()
	m.Open()

	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 1
	pm.TopicName = "/a/b/c"
	pm.Payload = []byte{0xBE, 0xEF, 0xED}
	pm.MessageID = 91
	key := ibound_mid2key(pm.MessageID)
	m.Put(key, pm)

	if len(m.messages) != 1 {
		t.Fatalf("message not in store")
	}
}

func Test_MemoryStore_Get(t *testing.T) {
	m := NewMemoryStore()
	m.Open()
	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 1
	pm.TopicName = "/a/b/c"
	pm.Payload = []byte{0xBE, 0xEF, 0xED}
	pm.MessageID = 120

	key := obound_mid2key(pm.MessageID)
	m.Put(key, pm)

	if len(m.messages) != 1 {
		t.Fatalf("message not in store")
	}

	exp := []byte{
		/* msg type */
		0x32, // qos 1

		/* remlen */
		0x0d,

		/* topic, msg id in varheader */
		0x00, // length of topic
		0x06,
		0x2F, // /
		0x61, // a
		0x2F, // /
		0x62, // b
		0x2F, // /
		0x63, // c

		/* msg id (is always 2 bytes) */
		0x00,
		0x78,

		/*payload */
		0xBE,
		0xEF,
		0xED,
	}

	msg := m.Get(key)

	if msg == nil {
		t.Fatalf("message not retreived from store")
	}

	var buf bytes.Buffer
	msg.Write(&buf)
	if !bytes.Equal(exp, buf.Bytes()) {
		t.Fatalf("message from store not same as what went in")
	}
}

func Test_MemoryStore_Del(t *testing.T) {
	m := NewMemoryStore()
	m.Open()

	pm := NewControlPacket(PUBLISH).(*PublishPacket)
	pm.Qos = 1
	pm.TopicName = "/a/b/c"
	pm.Payload = []byte{0xBE, 0xEF, 0xED}
	pm.MessageID = 17

	key := obound_mid2key(pm.MessageID)

	m.Put(key, pm)

	if len(m.messages) != 1 {
		t.Fatalf("message not in store")
	}

	m.Del(key)

	if len(m.messages) != 1 {
		t.Fatalf("message still exists after deletion")
	}
}
