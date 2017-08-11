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
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

func keepalive(c *client) {
	defer c.workers.Done()
	DEBUG.Println(PNG, "keepalive starting")
	var checkInterval time.Duration
	var pingSent time.Time

	if c.options.KeepAlive > 10*time.Second {
		checkInterval = 5 * time.Second
	} else {
		checkInterval = c.options.KeepAlive / 2
	}

	intervalTicker := time.NewTicker(checkInterval)
	defer intervalTicker.Stop()

	for {
		select {
		case <-c.stop:
			DEBUG.Println(PNG, "keepalive stopped")
			return
		case <-intervalTicker.C:
			if time.Now().Sub(c.lastSent) >= c.options.KeepAlive || time.Now().Sub(c.lastReceived) >= c.options.KeepAlive {
				if !c.pingOutstanding {
					DEBUG.Println(PNG, "keepalive sending ping")
					ping := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)
					//We don't want to wait behind large messages being sent, the Write call
					//will block until it it able to send the packet.
					c.pingOutstanding = true
					ping.Write(c.conn)
					c.lastSent = time.Now()
					pingSent = time.Now()
				}
			}
			if c.pingOutstanding && time.Now().Sub(pingSent) >= c.options.PingTimeout {
				CRITICAL.Println(PNG, "pingresp not received, disconnecting")
				c.errors <- errors.New("pingresp not received, disconnecting")
				return
			}
		}
	}
}
