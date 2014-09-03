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
	. "github.com/alsm/hrotti/packets"
)

type Message interface {
	Duplicate() bool
	Qos() byte
	Retained() bool
	Topic() string
	MessageID() uint16
	Payload() []byte
}

type message struct {
	duplicate bool
	qos       byte
	retained  bool
	topic     string
	messageID uint16
	payload   []byte
}

func (m *message) Duplicate() bool {
	return m.duplicate
}

func (m *message) Qos() byte {
	return m.qos
}

func (m *message) Retained() bool {
	return m.retained
}

func (m *message) Topic() string {
	return m.topic
}

func (m *message) MessageID() uint16 {
	return m.messageID
}

func (m *message) Payload() []byte {
	return m.payload
}

func messageFromPublish(p *PublishPacket) Message {
	return &message{
		duplicate: p.Dup,
		qos:       p.Qos,
		retained:  p.Retain,
		topic:     p.TopicName,
		messageID: p.MessageID,
		payload:   p.Payload,
	}
}

func newConnectMsgFromOptions(options ClientOptions) *ConnectPacket {
	//m := newMsg(CONNECT, false, 0, false)
	m := NewControlPacket(CONNECT).(*ConnectPacket)

	m.CleanSession = options.cleanSession
	m.WillFlag = options.willEnabled
	m.WillRetain = options.willRetained
	m.ClientIdentifier = options.clientId

	if options.willEnabled {
		m.WillQos = options.willQos
		m.WillTopic = options.willTopic
		m.WillMessage = options.willPayload
	}

	if options.username != "" {
		m.Username = options.username
		//mustn't have password without user as well
		if options.password != "" {
			m.Password = []byte(options.password)
		}
	}

	m.KeepaliveTimer = uint16(options.keepAlive)

	return m
}
