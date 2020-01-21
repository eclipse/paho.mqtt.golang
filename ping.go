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

	"github.com/desertbit/timer"
)

func keepalive(c *client) {
	defer c.workers.Done()
	DEBUG.Println(PNG, "keepalive starting")

	c.keepaliveTimer = timer.NewTimer(time.Duration(c.options.KeepAlive * int64(time.Second)))
	defer c.keepaliveTimer.Stop()

	c.pingTimeoutTimer = timer.NewTimer(c.options.PingTimeout)
	c.pingTimeoutTimer.Stop()
	defer c.pingTimeoutTimer.Stop()

	for {
		select {
		case <-c.stop:
			DEBUG.Println(PNG, "keepalive stopped")
			return
		case <-c.keepaliveTimer.C:
			resetKeepaliveTimer(c)
			DEBUG.Println(PNG, "keepalive sending ping")
			ping := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)
			//We don't want to wait behind large messages being sent, the Write call
			//will block until it it able to send the packet.
			ping.Write(c.conn)
			c.pingTimeoutTimer.Reset(c.options.PingTimeout)
		case <-c.pingTimeoutTimer.C:
			CRITICAL.Println(PNG, "pingresp not received, disconnecting")
			c.errors <- errors.New("pingresp not received, disconnecting")
			return
		}
	}
}

func resetKeepaliveTimer(c *client) {
	c.keepaliveTimer.Reset(time.Duration(c.options.KeepAlive * int64(time.Second)))
}

func stopPingTimeOutTimer(c *client) {
	c.pingTimeoutTimer.Stop()
}
