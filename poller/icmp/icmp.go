package icmp

import (
	"encoding/json"
	"errors"
	"log"
	M "probe/model"
	"probe/poller"
	"sync"
	"time"
)

var (
	Ping    = ping
	ErrICMP = errors.New("NO ICMP CONFIGURATION")
)

type Poller struct {
	in      M.Receiver[<-chan []byte]
	out     []M.Sender[any]
	workers int
	timeout time.Duration
}

func NewPoller(c M.Receiver[<-chan []byte], p []M.Sender[any], opts ...poller.OptFunc) M.Poller {
	return newICMPPoller(c, p, opts...)
}

func NewScanner(c M.Receiver[<-chan []byte], p []M.Sender[any], opts ...poller.OptFunc) M.Scanner {
	return newICMPPoller(c, p, opts...)
}

func newICMPPoller(c M.Receiver[<-chan []byte], p []M.Sender[any], opts ...poller.OptFunc) *Poller {
	options := poller.DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	return &Poller{in: c, out: p, workers: options.Workers, timeout: options.Timeout}
}

func (p *Poller) Poll() error {
	var wg sync.WaitGroup
	c := make(chan M.None, p.workers)
	var err error
	for data := range p.in.Receive() {
		var icmp M.ICMPReq
		err = json.Unmarshal(data, &icmp)
		if err != nil {
			for _, sender := range p.out {
				sender.Send(err)
			}
			continue
		}
		wg.Add(1)
		c <- M.None{}
		go func(i M.ICMPReq) {
			defer func() { wg.Done(); <-c }()
			if i.Params.Timeout == 0 {
				i.Params.Timeout = p.timeout
			}
			out, err := Ping(i)
			if err != nil {
				out.Err = err.Error()
			}
			log.Println("ICMP output -- ", out)
			// p.out.Send(out)
			for _, sender := range p.out {
				sender.Send(out)
			}
		}(icmp)
		log.Println("ICMP -- ", icmp)
	}
	wg.Wait()
	return nil
}

func (p *Poller) Scan(c any) (o any, err error) {
	switch v := c.(type) {
	case M.ICMPReq:
		o, err = Ping(v)
	default:
		err = ErrICMP
	}
	return
}
