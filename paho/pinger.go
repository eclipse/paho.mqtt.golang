package paho

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

type PingFailHandler func(error)

type Pinger interface {
	Start(net.Conn, time.Duration)
	Stop()
	PingResp()
}

type pingHandler struct {
	stop            chan struct{}
	conn            net.Conn
	lastPing        time.Time
	pingOutstanding int32
	pingFailHandler PingFailHandler
}

func DefaultPingerWithCustomFailHandler(pfh PingFailHandler) pingHandler {
	return pingHandler{pingFailHandler: pfh}
}

func PFH(err error) {
	log.Fatalln(err)
}

func (p *pingHandler) Start(c net.Conn, pt time.Duration) {
	debug.Println("pingHandler started")
	if p.pingFailHandler == nil {
		p.pingFailHandler = PFH
	}
	p.conn = c
	p.stop = make(chan struct{})
	checkTicker := time.NewTicker(pt / 4)
	defer checkTicker.Stop()
	for {
		select {
		case <-p.stop:
			debug.Println("pingHandler stopped")
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
					debug.Println("pingHandler sending ping request")
					p.pingFailHandler(err)
					return
				}
			}
		}
	}
}

func (p *pingHandler) Stop() {
	debug.Println("pingHandler stopping")
	select {
	case <-p.stop:
		//Already stopped, do nothing
	default:
		close(p.stop)
	}
}

func (p *pingHandler) PingResp() {
	debug.Println("pingHandler resetting pingOutstanding")
	atomic.StoreInt32(&p.pingOutstanding, 0)
}
