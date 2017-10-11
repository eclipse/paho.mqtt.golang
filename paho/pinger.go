package paho

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

type Pinger interface {
	Start()
	Stop()
	PingResp()
}

type PingHandler struct {
	stop            chan struct{}
	conn            net.Conn
	pingTimer       time.Duration
	lastPing        time.Time
	pingOutstanding int32
	pingFailHandler func(error)
}

func PFH(err error) {
	log.Fatalln(err)
}

func NewPingHandler(c net.Conn, pt time.Duration, pfh func(error)) Pinger {
	return &PingHandler{
		conn:            c,
		pingTimer:       pt,
		pingFailHandler: pfh,
	}
}

func (p *PingHandler) Start() {
	p.stop = make(chan struct{})
	checkTicker := time.NewTicker(p.pingTimer / 4)
	defer checkTicker.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-checkTicker.C:
			if atomic.LoadInt32(&p.pingOutstanding) > 0 && time.Now().Sub(p.lastPing) > p.pingTimer {
				p.pingFailHandler(fmt.Errorf("Ping resp timed out"))
				//ping outstanding and not reset in 1.5 times ping timer
				return
			}
			if time.Now().Sub(p.lastPing) >= p.pingTimer {
				//time to send a ping
				if err := packets.NewControlPacket(packets.PINGREQ).Send(p.conn); err != nil {
					p.pingFailHandler(err)
					return
				}
			}
		}
	}
}

func (p *PingHandler) Stop() {
	close(p.stop)
}

func (p *PingHandler) PingResp() {
	atomic.StoreInt32(&p.pingOutstanding, 0)
}
