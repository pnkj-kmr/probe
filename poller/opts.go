package poller

import (
	"runtime"
	"time"
)

type Opts struct {
	Workers int
	Timeout time.Duration
}

type OptFunc func(*Opts)

// DefaultOpts returns default options.
func DefaultOpts() Opts {
	return Opts{
		Workers: runtime.NumCPU() * 1000,
		Timeout: time.Second * 10,
	}
}

func WithWorkers(workers int) OptFunc {
	return func(opts *Opts) {
		opts.Workers = workers
	}
}

func WithTimeout(timeout time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.Timeout = timeout
	}
}
