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
	Wait()
	WaitTimeout(time.Duration)
	flowComplete()
}

type baseToken struct {
	m        sync.RWMutex
	complete chan struct{}
	ready    bool
}

// Wait will wait indefinitely for the Token to complete, ie the Publish
// to be sent and confirmed receipt from the broker
func (b *baseToken) Wait() {
	b.m.Lock()
	defer b.m.Unlock()
	if !b.ready {
		<-b.complete
		b.ready = true
	}
}

// WaitTimeout takes a time in ms
func (b *baseToken) WaitTimeout(d time.Duration) {
	b.m.Lock()
	defer b.m.Unlock()
	if !b.ready {
		select {
		case <-b.complete:
			b.ready = true
		case <-time.After(d):
		}
	}
}

func (b *baseToken) flowComplete() {
	close(b.complete)
}

func newToken(tType byte) Token {
	switch tType {
	case SUBSCRIBE:
		return &SubscribeToken{baseToken: baseToken{complete: make(chan struct{})}, subResult: make(map[string]byte)}
	case PUBLISH:
		return &PublishToken{baseToken: baseToken{complete: make(chan struct{})}}
	}
	return nil
}

type PublishToken struct {
	baseToken
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
