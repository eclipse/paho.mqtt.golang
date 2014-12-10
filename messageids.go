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
	"errors"
	"sync"
)

// MId is 16 bit message id as specified by the MQTT spec.
// In general, these values should not be depended upon by
// the client application.
type MId uint16

type messageIds struct {
	sync.Mutex
	idChan   chan MId
	index    map[MId]bool
	stopChan chan struct{}
}

const (
	MId_MAX MId = 65535
	MId_MIN MId = 1
)

func (mids *messageIds) stop() {
	close(mids.stopChan)
}

func (mids *messageIds) generateMsgIds() {

	mids.idChan = make(chan MId, 10)
	mids.stopChan = make(chan struct{})
	go func(mid *messageIds) {
		for {
			mid.Lock()
			for i := MId_MIN; i < MId_MAX; i++ {
				if !mid.index[i] {
					mid.index[i] = true
					mid.Unlock()
					select {
					case mid.idChan <- i:
					case <-mid.stopChan:
						return
					}
					break
				}
			}
		}
	}(mids)
}

func (mids *messageIds) freeId(id MId) {
	mids.Lock()
	defer mids.Unlock()
	//trace_v(MID, "freeing message id: %v", id)
	mids.index[id] = false
}

func (mids *messageIds) getId() (MId, error) {
	select {
	case i := <-mids.idChan:
		return i, nil
	case <-mids.stopChan:
	}
	return 0, errors.New("Failed to get next message id.")
}
