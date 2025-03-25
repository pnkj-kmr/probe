package M

import (
	"iter"
)

type StreamType interface {
	iter.Seq[[]byte] |
		<-chan []byte |
		chan<- []byte |
		[]byte
}

// <-chan any | iter.Seq[any]

type Poller interface {
	Poll() error
}

type Scanner interface {
	Scan(any) (any, error)
}

type Receiver[T StreamType] interface {
	Receive() T
}

type Sender[T StreamType] interface {
	Send(T) error
}
