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
	"fmt"
	"strconv"
)

const (
	_IBOUND_PRE = "i."
	_OBOUND_PRE = "o."
)

// Store is an interface which can be used to provide implementations
// for message persistence.
// Because we may have to store distinct messages with the same
// message ID, we need a unique key for each message. This is
// possible by prepending "i." or "o." to each message id
type Store interface {
	Open()
	Put(string, *Message)
	Get(string) *Message
	All() []string
	Del(string)
	Close()
	Reset()
}

// A key MUST have the form "X.[messageid]"
// where X is 'i' or 'o'
func key2mid(key string) MId {
	s := key[2:]
	i, err := strconv.Atoi(s)
	chkerr(err)
	return MId(i)
}

// Return a string of the form "i.[id]"
func ibound_mid2key(id MId) string {
	return fmt.Sprintf("%s%d", _IBOUND_PRE, id)
}

// Return a string of the form "o.[id]"
func obound_mid2key(id MId) string {
	return fmt.Sprintf("%s%d", _OBOUND_PRE, id)
}

// govern which outgoing messages are persisted
func persist_obound(s Store, m *Message) {
	switch m.QoS() {
	case QOS_ZERO:
		switch m.msgType() {
		case PUBACK:
			// Sending puback. delete matching publish
			// from ibound
			s.Del(ibound_mid2key(m.MsgId()))
		case PUBCOMP:
			// Sending pubcomp. delete matching pubrel
			// from ibound
			s.Del(ibound_mid2key(m.MsgId()))
		}
	case QOS_ONE:
		switch m.msgType() {
		case PUBLISH:
			// Sending publish. store in obound
			// until puback received
			s.Put(obound_mid2key(m.MsgId()), m)
		case PUBREL:
			// Sending pubrel. overwrite publish
			// in obound until pubcomp received
			s.Put(obound_mid2key(m.MsgId()), m)
		case SUBSCRIBE:
			// Sending subscribe. store in obound
			// until suback received
			s.Put(obound_mid2key(m.MsgId()), m)
		case UNSUBSCRIBE:
			// Sending unsubscribe. store in obound
			// until unsuback received
			s.Put(obound_mid2key(m.MsgId()), m)
		default:
			chkcond(false)
		}
	case QOS_TWO:
		switch m.msgType() {
		case PUBLISH:
			// Sending publish. store in obound
			// until pubrel received
			s.Put(obound_mid2key(m.MsgId()), m)
		default:
			chkcond(false)
		}
	}
}

// govern which incoming messages are persisted
func persist_ibound(s Store, m *Message) {
	switch m.QoS() {
	case QOS_ZERO:
		switch m.msgType() {
		case PUBACK:
			// Received a puback. delete matching publish
			// from obound
			s.Del(obound_mid2key(m.MsgId()))
		case SUBACK:
			// Received a suback. delete matching subscribe
			// from obound
			s.Del(obound_mid2key(m.MsgId()))
		case UNSUBACK:
			// Received a unsuback. delete matching unsubscribe
			// from obound
			s.Del(obound_mid2key(m.MsgId()))
		case PUBLISH:
		case PUBREC:
		case PUBCOMP:
			// Received a pubcomp. delete matching pubrel
			// from obound
			s.Del(obound_mid2key(m.MsgId()))
		case PINGRESP:
		case CONNACK:
		default:
			chkcond(false)
		}
	case QOS_ONE:
		switch m.msgType() {
		case PUBLISH:
			// Received a publish. store it in ibound
			// until puback sent
			s.Put(ibound_mid2key(m.MsgId()), m)
		case PUBREL:
			// Received a pubrel. Overwrite publish in ibound
			// until pubcomp sent
			s.Put(ibound_mid2key(m.MsgId()), m)
		default:
			chkcond(false)
		}
	case QOS_TWO:
		switch m.msgType() {
		case PUBLISH:
			// Received a publish. store it in ibound
			// until pubrel received
			s.Put(ibound_mid2key(m.MsgId()), m)
		default:
			chkcond(false)
		}
	}
}
