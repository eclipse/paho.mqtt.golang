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

// MId is 16 bit message id as specified by the MQTT spec.
// In general, these values should not be depended upon by
// the client application.
type MId uint16

type messageIds struct {
	sync.Mutex
	idChan chan MId
	index  map[MId]bool
}

const (
	MId_MAX MId = 65535
	MId_MIN MId = 1
)

func (mids *messageIds) generateMsgIds() {
	mids.idChan = make(chan MId, 10)
	go func() {
		for {
			mids.Lock()
			for i := MId_MIN; i < MId_MAX; i++ {
				if !mids.index[i] {
					mids.index[i] = true
					mids.Unlock()
					mids.idChan <- i
					break
				}
			}
		}
	}()
}

func (mids *messageIds) freeId(id MId) {
	mids.Lock()
	defer mids.Unlock()
	//trace_v(MID, "freeing message id: %v", id)
	mids.index[id] = false
}

func (mids *messageIds) getId() MId {
	return <-mids.idChan
}
