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
	p.ctx.WithEngine(p.enigne)
	log.Println("setting... context")

	// DB engine init
	err = newDBEngine(p.ctx, options)
	if err != nil {
		return nil, err
	}
	log.Println("setting... db")

	//Export engine init
	// workers node
	err = newExportEngine(p.ctx, options)
	if err != nil {
		return nil, err
	}
	log.Println("setting... export")

	// Poller engine init
	err = newPollEngine(p.ctx, options)
	if err != nil {
		return nil, err
	}
	log.Println("setting... poll")

	// Schedule engine init
	err = newScheduleEngine(p.ctx, options)
	if err != nil {
		return nil, err
	}
	log.Println("setting... schedule")

	// API engine init
	err = newAPIEngine(p.ctx, options)
	if err != nil {
		return nil, err
	}
	log.Println("setting... api")

	// Event engine init
	err = newEventEngine(p.ctx, options)
	if err != nil {
		return nil, err
	}
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

	<-p.ctx.Context().Done()
}

func (p *Probe) Stop() {
	log.Println("[PROBE] Shutting down gracefully...")

	p.ctx.Schedule().Stop()
	p.ctx.Export().Stop()
	p.ctx.API().Stop()

	log.Println("[PROBE] gracefully shutdown")

}
