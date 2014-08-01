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
	"time"
)

type lastcontact struct {
	sync.Mutex
	lasttime time.Time
}

func (l *lastcontact) update() {
	l.Lock()
	defer l.Unlock()
	l.lasttime = time.Now()

}

func (l *lastcontact) get() time.Time {
	l.Lock()
	defer l.Unlock()
	return l.lasttime
}

func newPingReqMsg() *Message {
	m := newMsg(PINGREQ, false, QOS_ZERO, false)
	m.remlen = uint32(0)
	return m
}

func keepalive(c *MqttClient) {
	DEBUG.Println(PNG, "keepalive starting")

	for {
		select {
		case <-c.stop:
			WARN.Println(PNG, "keepalive stopped")
			return
		default:
			last := uint(time.Since(c.lastContact.get()).Seconds())
			//DEBUG.Println(PNG, "last contact: %d (timeout: %d)", last, c.options.timeout)
			if last > c.options.keepAlive {
				if !c.pingOutstanding {
					DEBUG.Println(PNG, "keepalive sending ping")
					ping := newPingReqMsg()
					c.oboundP <- ping
					c.pingOutstanding = true
				} else {
					CRITICAL.Println(PNG, "pingresp not received, disconnecting")
					go c.options.onconnlost(c, errors.New("pingresp not received, disconnecting"))
					c.disconnect()
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}
