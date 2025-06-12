package probe

import (
	"context"
)

type Opts struct {
	context context.Context
	port    int
	api     bool
	icmp    bool
	snmp    bool
	// interval    time.Duration
	maxRestarts     int
	workers         int
	kafka           bool
	maxBucketSize   int
	totalPartitions int
}

type OptFunc func(*Opts)

// DefaultOpts returns default options.
func DefaultOpts() Opts {
	// TODO need to find default options from db
	// env.db will hold all options with default values
	return Opts{
		context: context.Background(),
		port:    8080,
		api:     true,
		kafka:   true,
		// interval:    time.Second * 10,
		maxRestarts: 3,
		workers:     100,
		// TODO - need to take from app.yml
		maxBucketSize:   2,
		totalPartitions: 2,
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

func WithKafka() OptFunc {
	return func(opts *Opts) {
		opts.kafka = true
	}
}

func WithBucketSize(t int) OptFunc {
	return func(opts *Opts) {
		opts.maxBucketSize = t
	}
}

func WithSNMP() OptFunc {
	return func(opts *Opts) {
		opts.snmp = true
	}
}

// func WithInterval(t time.Duration) OptFunc {
// 	return func(opts *Opts) {
// 		opts.interval = t
// 	}
// }

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
