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
	in      M.Consumer[<-chan []byte]
	out     M.Producer[[]byte]
	workers int
	timeout time.Duration
}

func NewPoller(c M.Consumer[<-chan []byte], p M.Producer[[]byte], opts ...poller.OptFunc) M.Poller {
	return newICMP(c, p, opts...)
}

func NewScanner(c M.Consumer[<-chan []byte], p M.Producer[[]byte], opts ...poller.OptFunc) M.Scanner {
	return newICMP(c, p, opts...)
}

func newICMP(c M.Consumer[<-chan []byte], p M.Producer[[]byte], opts ...poller.OptFunc) *Poller {
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
	for data := range p.in.Consume() {
		var icmp M.InICMP
		err = json.Unmarshal(data, &icmp)
		if err != nil {
			d, _ := json.Marshal(M.OutICMP{Err: err.Error()})
			p.out.Produce(d)
		} else {
			wg.Add(1)
			c <- M.None{}
			go func(i M.InICMP) {
				defer func() { wg.Done(); <-c }()
				if i.Timeout == 0 {
					i.Timeout = p.timeout
				}
				if i.Count == 0 {
					i.Count = 4
				}
				out, err := Ping(i)
				if err != nil {
					out.Err = err.Error()
				}
				log.Println("ICMP output -- ", out)
				d, _ := json.Marshal(out)
				p.out.Produce(d)
			}(icmp)
			log.Println("ICMP -- ", icmp)
		}

	}
	wg.Wait()
	return nil
}

func (p *Poller) Scan(c any) (o any, err error) {
	switch v := c.(type) {
	case M.InICMP:
		o, err = Ping(v)
	default:
		err = ErrICMP
	}
	return
}
