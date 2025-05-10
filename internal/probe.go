package probe

import (
	"context"
	"log"

	"github.com/anthdm/hollywood/actor"
)

// Probe
type Probe struct {
	ctx    *Context
	enigne *actor.Engine
}

func NewProbeWithContext(c context.Context, opts ...OptFunc) (*Probe, error) {
	opts = append(opts, WithContext(c))
	return NewProbe(opts...)
}

func NewProbe(opts ...OptFunc) (*Probe, error) {
	p := &Probe{}

	e, err := newEngine()
	if err != nil {
		return p, nil
	}
	p.enigne = e
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	p.ctx = newContext(options.context)
	log.Println("setting... context")

	// DB engine init
	db, err := newDBEngine(p.enigne, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithDB(db)
	log.Println("setting... db")

	//Export engine init
	// workers node
	export, err := newExportEngine(p.enigne, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithExport(export)
	log.Println("setting... export")

	// Poller engine init
	poll, err := newPollEngine(db, export, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithPoll(poll)
	log.Println("setting... poll")

	// Schedule engine init
	schedule, err := newScheduleEngine(p.enigne, poll, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithSchdule(schedule)
	log.Println("setting... schedule")

	// API engine init
	api, err := newAPIEngine(p.enigne, db, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithAPI(api)
	log.Println("setting... api")

	// Event engine init
	event, err := newEventEngine(p.enigne, opts...)
	if err != nil {
		return nil, err
	}
	p.ctx.WithEvent(event)
	log.Println("setting... event")

	return p, nil
}

func (p *Probe) Context() *Context {
	return p.ctx
}

func (p *Probe) Start() {
	log.Println("[PROBE] started")

	p.ctx.Schedule().Start()
	p.ctx.Export().Start()
	p.ctx.API().Start()
	p.ctx.DB().Start()

	<-p.ctx.context.Done()
}

func (p *Probe) Stop() {
	log.Println("[PROBE] Shutting down gracefully...")

	p.ctx.Schedule().Stop()
	p.ctx.Export().Stop()
	p.ctx.API().Stop()
	p.ctx.DB().Stop()

	log.Println("[PROBE] gracefully shutdown")

}
