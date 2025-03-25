package probe

import (
	"context"
	"time"
)

type Opts struct {
	context     context.Context
	port        int
	api         bool
	icmp        bool
	snmp        bool
	interval    time.Duration
	maxRestarts int
	workers     int
}

type OptFunc func(*Opts)

// DefaultOpts returns default options.
func DefaultOpts() Opts {
	return Opts{
		context:     context.Background(),
		port:        8080,
		api:         true,
		interval:    time.Second * 10,
		maxRestarts: 3,
		workers:     100,
	}
}

func WithContext(ctx context.Context) OptFunc {
	return func(opts *Opts) {
		opts.context = ctx
	}
}

func WithICMP() OptFunc {
	return func(opts *Opts) {
		opts.icmp = true
	}
}

func WithSNMP() OptFunc {
	return func(opts *Opts) {
		opts.snmp = true
	}
}

func WithInterval(t time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.interval = t
	}
}

func WithMaxRestarts(t int) OptFunc {
	return func(opts *Opts) {
		opts.maxRestarts = t
	}
}

func WithWorkers(t int) OptFunc {
	return func(opts *Opts) {
		opts.workers = t
	}
}
