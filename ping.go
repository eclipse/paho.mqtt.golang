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
	DEBUG.Println(PNG, "keepalive starting")

	pingTimer := timer{Timer: time.NewTimer(c.options.KeepAlive)}
	pingRespTimer := timer{Timer: time.NewTimer(c.options.PingTimeout)}

	pingRespTimer.Stop()

	for {
		select {
		case <-c.stop:
			DEBUG.Println(PNG, "keepalive stopped")
			c.workers.Done()
			return
		case <-pingTimer.C:
			pingTimer.SetRead(true)
			DEBUG.Println(PNG, "keepalive sending ping")
			ping := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)
			//We don't want to wait behind large messages being sent, the Write call
			//will block until it it able to send the packet.
			ping.Write(c.conn)

			pingRespTimer.Reset(c.options.PingTimeout)
		case <-c.pingResp:
			DEBUG.Println(NET, "resetting ping timers")
			pingRespTimer.Stop()
			pingTimer.Reset(c.options.KeepAlive)
		case <-pingRespTimer.C:
			pingRespTimer.SetRead(true)
			CRITICAL.Println(PNG, "pingresp not received, disconnecting")
			c.workers.Done()
			c.internalConnLost(errors.New("pingresp not received, disconnecting"))
			pingTimer.Stop()
			return
		}
	}
}

type timer struct {
	*time.Timer
	readFrom bool
}

func (t *timer) SetRead(v bool) {
	t.readFrom = v
}

func (t *timer) Stop() bool {
	defer t.SetRead(true)

	if !t.Timer.Stop() && !t.readFrom {
		<-t.C
		return false
	}
	return true
}

func (t *timer) Reset(d time.Duration) bool {
	defer t.SetRead(false)
	t.Stop()
	return t.Timer.Reset(d)
}
