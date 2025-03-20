package M

import (
	"iter"
)

type StreamType interface {
	<-chan any | <-chan []byte | iter.Seq[any] | iter.Seq[[]byte]
}

type Poller interface {
	Poll() error
}

type Scanner interface {
	Scan(any) (any, error)
}

type Consumer[T StreamType] interface {
	Consume() T
}

type Producer interface {
	Produce(any) error
}
