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
	"sync"
	"time"
)

// MId is 16 bit message id as specified by the MQTT spec.
// In general, these values should not be depended upon by
// the client application.
type MId uint16

type messageIds struct {
	sync.RWMutex
	index [65535]Token
}

func (mids *messageIds) cleanUp() {
	mids.Lock()
	for _, token := range mids.index {
		switch t := token.(type) {
		case *PublishToken:
			t.err = fmt.Errorf("Connection lost before Publish completed")
		case *SubscribeToken:
			t.err = fmt.Errorf("Connection lost before Subscribe completed")
		case *UnsubscribeToken:
			t.err = fmt.Errorf("Connection lost before Unsubscribe completed")
		case nil:
			continue
		}
		token.flowComplete()
	}
	mids.index = [65535]Token{}
	mids.Unlock()
}

func (mids *messageIds) freeID(id uint16) {
	mids.Lock()
	mids.index[id-1] = nil
	mids.Unlock()
}

func (mids *messageIds) getID(t Token) uint16 {
	mids.Lock()
	defer mids.Unlock()
	for i, s := range mids.index {
		if s == nil {
			mids.index[i] = t
			return uint16(i + 1)
		}
	}
	return 0
}

func (mids *messageIds) getToken(id uint16) Token {
	mids.RLock()
	defer mids.RUnlock()
	t := mids.index[id-1]
	if t == nil {
		t = &DummyToken{id: id}
	}
	return t
}

type DummyToken struct {
	id uint16
}

func (d *DummyToken) Wait() bool {
	return true
}

func (d *DummyToken) WaitTimeout(t time.Duration) bool {
	return true
}

func (d *DummyToken) flowComplete() {
	ERROR.Printf("A lookup for token %d returned nil\n", d.id)
}

func (d *DummyToken) Error() error {
	return nil
}
