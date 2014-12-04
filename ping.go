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
	. "github.com/alsm/hrotti/packets"
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

func keepalive(c *MqttClient) {
	defer c.workers.Done()
	DEBUG.Println(PNG, "keepalive starting")

	for {
		select {
		case <-c.stop:
			DEBUG.Println(PNG, "keepalive stopped")
			return
		default:
			last := uint(time.Since(c.lastContact.get()).Seconds())
			//DEBUG.Printf("%s last contact: %d (timeout: %d)", PNG, last, uint(c.options.KeepAlive.Seconds()))
			if last > uint(c.options.KeepAlive.Seconds()) {
				if !c.pingOutstanding {
					DEBUG.Println(PNG, "keepalive sending ping")
					ping := NewControlPacket(PINGREQ).(*PingreqPacket)
					//We don't want to wait behind large messages being sent, the Write call
					//will block until it it able to send the packet.
					ping.Write(c.conn)
					c.pingOutstanding = true
				} else {
					CRITICAL.Println(PNG, "pingresp not received, disconnecting")
					go c.options.OnConnectionLost(c, errors.New("pingresp not received, disconnecting"))
					go c.reconnect()
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}
