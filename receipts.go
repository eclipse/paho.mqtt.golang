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
	"sync"
)

// Receipt is a sort of token object that you will receive
// upon delivery of a published message.
type Receipt struct {
}

// A synchronized map for MId => chan Receipt
type receiptMap struct {
	sync.RWMutex
	contents map[MId]chan Receipt
}

func newReceiptMap() *receiptMap {
	return &receiptMap{contents: make(map[MId]chan Receipt, 128)}
}

func (m *receiptMap) put(mid MId, c chan Receipt) {
	m.Lock()
	defer m.Unlock()
	m.contents[mid] = c
}

func (m *receiptMap) get(mid MId) chan Receipt {
	m.RLock()
	defer m.RUnlock()
	return m.contents[mid]
}

func (m *receiptMap) end(mid MId) {
	m.Lock()
	defer m.Unlock()
	close(m.contents[mid])
	delete(m.contents, mid)
}
