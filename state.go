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
	"strconv"
)

// Look into the store and finish any flows found in there.
func (c *MqttClient) resume() []Receipt {
	DEBUG.Println(STA, "resuming client from stored state")

	keys := c.persist.All()
	DEBUG.Println(STA, "there are", len(keys), "persisted msgs found")

	out := []string{}
	in := []string{}
	for _, k := range keys {
		if k[0] == 'i' {
			in = append(in, k)
		} else {
			out = append(out, k)
		}
	}

	out = reorder(out)
	in = reorder(in)
	receipts := []Receipt{}

	// take care of inbound qos 2
	for i := 0; i < len(in); i++ {
		m := c.persist.Get(in[i])
		if m.QoS() == QOS_TWO {
			DEBUG.Println(STA, "resume inbound qos2")
			c.ibound <- m
			// there needs to be some identifying information for receipts
			receipts = append(receipts, Receipt{})
		}
	}

	// take care of inbound qos 1
	for i := 0; i < len(in); i++ {
		m := c.persist.Get(in[i])
		if m.QoS() == QOS_ONE {
			DEBUG.Println(STA, "resume inbound qos1")
			c.ibound <- m
			// there needs to be some identifying information for receipts
			receipts = append(receipts, Receipt{})
		}
	}

	// take care of outbound qos 2
	for i := 0; i < len(out); i++ {
		m := c.persist.Get(out[i])
		if m.QoS() == QOS_TWO {
			DEBUG.Println(STA, "resume outbound qos2")
			if c.receipts.get(m.MsgId()) == nil { // will be nil if client crashed
				c.receipts.put(m.MsgId(), make(chan Receipt, 1))
			}
			sendSendableWithTimeout(c.obound, sendable{m, c.receipts.get(m.MsgId())})
		}
	}

	// take care of outbound qos 1
	for i := 0; i < len(out); i++ {
		m := c.persist.Get(out[i])
		if m.QoS() == QOS_ONE {
			DEBUG.Println(STA, "resume outbound qos1")
			if c.receipts.get(m.MsgId()) == nil { // will be nil if client crashed
				c.receipts.put(m.MsgId(), make(chan Receipt, 1))
			}
			sendSendableWithTimeout(c.obound, sendable{m, c.receipts.get(m.MsgId())})
		}
	}

	DEBUG.Println(STA, "done resuming client")
	return receipts
}

// Reorder keys where the keys are reordered in
// STRICT NUMERICAL ORDER only.
// This is intended to be used as a helper to the reorder function,
// which will put the keys in proper MQTT order.
func sortKeys(keys []string) {
	// Bubble sort is good enough because we don't expect a lot
	// of keys.
	for i := 0; i < len(keys)-1; i++ {
		for j := 0; j < len(keys)-1; j++ {
			if !isAlessB(keys[j], keys[j+1]) {
				temp := keys[j+1]
				keys[j+1] = keys[j]
				keys[j] = temp
			}
		}
	}
}

// Takes two keys, and returns true if a < b,
// or false otherwise.
func isAlessB(a, b string) bool {
	// chop off leading "X."
	nA, _ := strconv.Atoi(a[2:])
	nB, _ := strconv.Atoi(b[2:])
	return nA < nB
}

// Create a new list of keys where the keys are reordered according
// to their message id's.
func reorder(keys []string) []string {
	reordered := []string{}
	if len(keys) == 0 {
		return reordered
	}

	sortKeys(keys) // put keys in order of msgids low to high

	prevId := MId(0)
	maxGap := 0
	maxGapIdx := 0

	for i := 0; i < len(keys); i++ {
		curId := key2mid(keys[i])
		gap := int(curId - prevId)
		if gap > maxGap {
			maxGap = gap
			maxGapIdx = i
		}
		prevId = curId
	}

	minId := key2mid(keys[0])
	maxId := prevId

	// check edge case where maximum gap is actually looped around
	if int(MId_MAX-maxId+minId) > maxGap {
		maxGapIdx = 0
	}

	// start adding messages
	for i := maxGapIdx; i < len(keys); i++ {
		reordered = append(reordered, keys[i])
	}

	// loop around, adding remaining messages
	for i := 0; i < maxGapIdx; i++ {
		reordered = append(reordered, keys[i])
	}

	return reordered
}
