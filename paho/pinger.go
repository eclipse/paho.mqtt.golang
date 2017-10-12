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
	Start(time.Duration)
	Stop()
	PingResp()
}

type PingHandler struct {
	stop            chan struct{}
	conn            net.Conn
	lastPing        time.Time
	pingOutstanding int32
	pingFailHandler func(error)
}

func PFH(err error) {
	log.Fatalln(err)
}

func NewPingHandler(c net.Conn, pfh func(error)) Pinger {
	return &PingHandler{
		conn:            c,
		pingFailHandler: pfh,
	}
}

func (p *PingHandler) Start(pt time.Duration) {
	p.stop = make(chan struct{})
	checkTicker := time.NewTicker(pt / 4)
	defer checkTicker.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-checkTicker.C:
			if atomic.LoadInt32(&p.pingOutstanding) > 0 && time.Now().Sub(p.lastPing) > pt {
				p.pingFailHandler(fmt.Errorf("Ping resp timed out"))
				//ping outstanding and not reset in 1.5 times ping timer
				return
			}
			if time.Now().Sub(p.lastPing) >= pt {
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
