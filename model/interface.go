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

type Consumer[T StreamType] interface {
	Consume() T
}

type Producer[T StreamType] interface {
	Produce(T) error
}

// type Receiver[T StreamType] interface {
// 	Receive() T
// }

type Exporter[T StreamType] interface {
	Export() T
}
