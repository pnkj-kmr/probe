package M

import (
	"iter"
)

type StreamType interface {
	iter.Seq[[]byte] |
		<-chan []byte | <-chan any |
		chan<- []byte | chan<- any |
		[]byte | any
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

type Exporter[T, V StreamType] interface {
	Spin()
	Push() T
	Export(V) error
	Close() error
}

type Finder interface {
	Find(string) ([]byte, error)
}

type DB interface {
	Finder

	Receive() <-chan []byte
	Count() uint64
	Create(string, []byte) error
	CreateBulk(map[string][]byte) (int, error)
	Update(string, []byte) error
	Delete(string) error
	DeleteAll() error
	FindAll() (map[string][]byte, error)
}
