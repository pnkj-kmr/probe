package snmp

import (
	"encoding/json"
	"errors"
	"log"
	M "probe/model"
	"probe/poller"
	"sync"
)

var (
	DoSNMP  = doSNMP
	ErrSNMP = errors.New("NO SNMP CONFIGURATION")
)

type Poller struct {
	in      M.Receiver[<-chan []byte]
	out     M.Sender[any]
	profile M.Finder
	workers int

	profileMap map[string]M.AuthSNMP
}

func NewPoller(c M.Receiver[<-chan []byte], p M.Sender[any], f M.Finder, opts ...poller.OptFunc) M.Poller {
	return newSNMPPoller(c, p, f, opts...)
}

func NewScanner(c M.Receiver[<-chan []byte], p M.Sender[any], f M.Finder, opts ...poller.OptFunc) M.Scanner {
	return newSNMPPoller(c, p, f, opts...)
}

func newSNMPPoller(c M.Receiver[<-chan []byte], p M.Sender[any], f M.Finder, opts ...poller.OptFunc) *Poller {
	options := poller.DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	profileMap := make(map[string]M.AuthSNMP)
	return &Poller{in: c, out: p, workers: options.Workers, profile: f, profileMap: profileMap}
}

func (p *Poller) Poll() error {
	var wg sync.WaitGroup
	c := make(chan M.None, p.workers)
	var err error
	for data := range p.in.Receive() {
		var snmp M.SNMPReq
		err = json.Unmarshal(data, &snmp)
		if err != nil {
			p.out.Send(err)
			continue
		}
		// getting profile
		lp, err := p.getFindProfile(snmp.LoginProfileID)
		if err != nil {
			p.out.Send(err)
			continue
		}
		wg.Add(1)
		c <- M.None{}
		go func(i M.SNMPReq, lp M.AuthSNMP) {
			defer func() { wg.Done(); <-c }()
			out, err := DoSNMP(i, lp)
			if err != nil {
				out.Err = err.Error()
			}
			log.Println("SNMP output -- ", out)
			p.out.Send(out)
		}(snmp, lp)
		log.Println("SNMP -- ", snmp)
	}
	wg.Wait()
	return nil
}

func (p *Poller) Scan(c any) (o any, err error) {
	switch v := c.(type) {
	case M.SNMPReq:
		lp, _err := p.getFindProfile(v.LoginProfileID)
		if _err != nil {
			err = _err
			return
		}
		o, err = DoSNMP(v, lp)
	default:
		err = ErrSNMP
	}
	return
}

func (p *Poller) getFindProfile(name string) (profile M.AuthSNMP, err error) {
	profile, ok := p.profileMap[name]
	if !ok {
		data, _err := p.profile.Find(name)
		if _err != nil {
			return profile, _err
		}
		_err = json.Unmarshal(data, &profile)
		if _err != nil {
			return profile, _err
		}
		p.profileMap[name] = profile
		return
	}
	return profile, nil
}
