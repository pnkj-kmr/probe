package probe

import (
	"context"
	"log"
)

// Probe
type Probe struct {
	ctx *Context
}

func NewProbeWithContext(c context.Context, opts ...OptFunc) (*Probe, error) {
	opts = append(opts, WithContext(c))
	return NewProbe(opts...)
}

func NewProbe(opts ...OptFunc) (*Probe, error) {
	p := &Probe{}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	p.ctx = newContext(options.context)

	// DB engine init
	db, err := newDBEngine(opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithDB(db)

	//Export engine init
	// workers node
	export, err := newExportEngine(8, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithExport(export)

	// Poller engine init
	poll, err := newPollEngine(db, export, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithPoll(poll)

	// Schedule engine init
	schedule, err := newScheduleEngine(poll, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithSchdule(schedule)

	// API engine init
	api, err := newAPIEngine(db, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithAPI(api)

	// Event engine init
	event, err := newEventEngine(opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithEvent(event)

	return p, nil
}

func (p *Probe) Context() *Context {
	return p.ctx
}

func (p *Probe) Start() {
	log.Println("[PROBE] started")

	// p.ctx.Schedule().Start()
	p.ctx.Export().Start()
	p.ctx.API().Start()
	p.ctx.DB().Start()

	<-p.ctx.context.Done()
}

func (p *Probe) Stop() {
	log.Println("[PROBE] Shutting down gracefully...")

	// p.ctx.Schedule().Stop()
	p.ctx.Export().Stop()
	p.ctx.API().Stop()
	p.ctx.DB().Stop()

	log.Println("[PROBE] gracefully shutdown")

}
