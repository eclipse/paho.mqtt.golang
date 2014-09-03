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
	. "github.com/alsm/hrotti/packets"
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
	Put(string, ControlPacket)
	Get(string) ControlPacket
	All() []string
	Del(string)
	Close()
	Reset()
}

// A key MUST have the form "X.[messageid]"
// where X is 'i' or 'o'
func key2mid(key string) uint16 {
	s := key[2:]
	i, err := strconv.Atoi(s)
	chkerr(err)
	return uint16(i)
}

// Return a string of the form "i.[id]"
func ibound_mid2key(id uint16) string {
	return fmt.Sprintf("%s%d", _IBOUND_PRE, id)
}

// Return a string of the form "o.[id]"
func obound_mid2key(id uint16) string {
	return fmt.Sprintf("%s%d", _OBOUND_PRE, id)
}

// govern which outgoing messages are persisted
func persist_obound(s Store, m ControlPacket) {
	switch m.Details().Qos {
	case 0:
		switch m.(type) {
		case *PubackPacket, *PubcompPacket:
			// Sending puback. delete matching publish
			// from ibound
			s.Del(ibound_mid2key(m.Details().MessageID))
		}
	case 1:
		switch m.(type) {
		case *PublishPacket, *PubrelPacket, *SubscribePacket, *UnsubscribePacket:
			// Sending publish. store in obound
			// until puback received
			s.Put(obound_mid2key(m.Details().MessageID), m)
		default:
			chkcond(false)
		}
	case 2:
		switch m.(type) {
		case *PublishPacket:
			// Sending publish. store in obound
			// until pubrel received
			s.Put(obound_mid2key(m.Details().MessageID), m)
		default:
			chkcond(false)
		}
	}
}

// govern which incoming messages are persisted
func persist_ibound(s Store, m ControlPacket) {
	switch m.Details().Qos {
	case 0:
		switch m.(type) {
		case *PubackPacket, *SubackPacket, *UnsubackPacket, *PubcompPacket:
			// Received a puback. delete matching publish
			// from obound
			s.Del(obound_mid2key(m.Details().MessageID))
		case *PublishPacket, *PubrecPacket, *PingrespPacket, *ConnackPacket:
		default:
			chkcond(false)
		}
	case 1:
		switch m.(type) {
		case *PublishPacket, *PubrelPacket:
			// Received a publish. store it in ibound
			// until puback sent
			s.Put(ibound_mid2key(m.Details().MessageID), m)
		default:
			chkcond(false)
		}
	case 2:
		switch m.(type) {
		case *PublishPacket:
			// Received a publish. store it in ibound
			// until pubrel received
			s.Put(ibound_mid2key(m.Details().MessageID), m)
		default:
			chkcond(false)
		}
	}
}
