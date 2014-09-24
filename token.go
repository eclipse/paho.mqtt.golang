/*
 * Copyright (c) 2014 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Allan Stockdill-Mander
 */

package mqtt

import (
	. "github.com/alsm/hrotti/packets"
	"sync"
	"time"
)

type PacketAndToken struct {
	p ControlPacket
	t Token
}

type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	flowComplete()
	Error() error
}

type baseToken struct {
	m        sync.RWMutex
	complete chan struct{}
	ready    bool
	err      error
}

// Wait will wait indefinitely for the Token to complete, ie the Publish
// to be sent and confirmed receipt from the broker
func (b *baseToken) Wait() bool {
	b.m.Lock()
	defer b.m.Unlock()
	if !b.ready {
		<-b.complete
		b.ready = true
	}
	return b.ready
}

// WaitTimeout takes a time in ms
func (b *baseToken) WaitTimeout(d time.Duration) bool {
	b.m.Lock()
	defer b.m.Unlock()
	if !b.ready {
		select {
		case <-b.complete:
			b.ready = true
		case <-time.After(d):
		}
	}
	return b.ready
}

func (b *baseToken) flowComplete() {
	close(b.complete)
}

func (b *baseToken) Error() error {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.err
}

func newToken(tType byte) Token {
	switch tType {
	case SUBSCRIBE:
		return &SubscribeToken{baseToken: baseToken{complete: make(chan struct{})}, subResult: make(map[string]byte)}
	case PUBLISH:
		return &PublishToken{baseToken: baseToken{complete: make(chan struct{})}}
	case UNSUBSCRIBE:
		return &UnsubscribeToken{baseToken: baseToken{complete: make(chan struct{})}}
	}
	return nil
}

type PublishToken struct {
	baseToken
	messageId uint16
}

func (p *PublishToken) MessageId() uint16 {
	return p.messageId
}

type SubscribeToken struct {
	baseToken
	subs      []string
	subResult map[string]byte
}

func (s *SubscribeToken) Result() map[string]byte {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.subResult
}

type UnsubscribeToken struct {
	baseToken
}
