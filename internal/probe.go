package probe

import (
	"context"
	"fmt"
)

// Probe
type Probe struct {
	ctx      *Context
	db       *DBEngine
	poll     *PollEngine
	schedule *ScheduleEngine
	export   *ExportEngine
}

// ProbeConfig holds the configuration of the probe.
type ProbeConfig struct {
	ctx context.Context
}

// NewProbeConfig returns a new default ProbeConfig.
func NewProbeConfig(ctx context.Context) ProbeConfig {
	if ctx == nil {
		ctx = context.Background()
	}
	return ProbeConfig{ctx: ctx}
}

func NewProbe(config ProbeConfig, opts ...OptFunc) (*Probe, error) {
	p := &Probe{}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	p.ctx = newContext(config.ctx)

	// DB engine init
	db, err := NewDBEngine(opts...)
	if err != nil {
		return nil, err
	}
	p.db = db

	//Export engine init
	// workers node
	export, err := NewExportEngine(8, opts...)
	if err != nil {
		return nil, err
	}
	p.export = export

	// Poller engine init
	poll, err := NewPollEngine(db, export, opts...)
	if err != nil {
		return nil, err
	}
	p.poll = poll

	// Schedule engine init
	schedule, err := NewScheduleEngine(poll, opts...)
	if err != nil {
		return nil, err
	}
	p.schedule = schedule
	return p, nil
}

func (p *Probe) Start() {
	// Simulating work by sleeping
	p.schedule.scheduler.ForEach(func(i int, scheduler *schedule) {
		fmt.Println("start----", i)
		scheduler.Start()
	})

	<-p.ctx.context.Done()
	fmt.Println("Shutting down gracefully...")
}

func (p *Probe) Stop() {
	fmt.Println("Shutting down gracefully...")
	p.schedule.scheduler.ForEach(func(i int, scheduler *schedule) {
		fmt.Println("stop----", i)
		scheduler.Stop()
	})
	fmt.Println("probe stop initiated...")
}
