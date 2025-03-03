package icmp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	M "probe/model"
	"sync"
	"time"
)

var (
	Ping   = ping
	NoICMP = errors.New("No ICMP configuration found")
)

type ctx struct {
	db       M.DBView
	bucket   string
	receiver M.Receiver
	workers  int
	timeout  time.Duration
}

func NewICMPPoller(db M.DBView, bucket string, receiver M.Receiver, workers int, timeout time.Duration) M.Poller {
	return &ctx{db, bucket, receiver, workers, timeout}
}

func NewICMPScanner() M.Scanner {
	return &ctx{}
}

func (p *ctx) Poll() (err error) {
	// getting configuration from db
	data, err := p.getConfigs()
	if err != nil {
		return err
	}

	// performing do poll for all stats
	err = p.doPoll(data)
	if err != nil {
		return err
	}
	return
}

func (p *ctx) getConfigs() (out []M.ICMPConf, err error) {
	data := p.db.Cursor(p.bucket)
	// TODO - need to validate the ICMPConf object
	var icmp M.ICMPConf
	for d := range data {
		err := json.Unmarshal(d, &icmp)
		if err != nil {
			log.Println("unable to convert as ICMPConf", err)
		}
		out = append(out, icmp)
	}
	if len(out) == 0 {
		return out, NoICMP
	}
	return
}

func (p *ctx) doPoll(data []M.ICMPConf) (err error) {
	var wg sync.WaitGroup
	c := make(chan M.None, p.workers)
	for i := 0; i < len(data); i++ {
		wg.Add(1)
		c <- M.None{}
		go func(i M.ICMPConf) {
			defer func() { wg.Done(); <-c }()
			if i.Timeout == 0 {
				i.Timeout = p.timeout
			}
			if i.Count == 0 {
				i.Count = 4 // default packet count if zero
			}
			out, err := Ping(i)
			if err != nil {
				out.Err = err.Error()
			}
			p.receiver.Receive() <- out
		}(data[i])
		log.Println("ICMP -- ", data[i], i+1)
	}
	wg.Wait()
	return nil
}

func (p *ctx) Scan(c interface{}) (o interface{}, err error) {
	// performing scan on given configuration
	switch v := c.(type) {
	case M.ICMPConf:
		o, err = Ping(v)
	default:
		err = errors.New(fmt.Sprint(v))
	}
	return
}
