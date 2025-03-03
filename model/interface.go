package M

import "io"

type Poller interface {
	Poll() error
}

type Scanner interface {
	Scan(interface{}) (interface{}, error)
}

type Receiver interface {
	Receive() chan<- interface{}
}

type Creator interface {
	Create(string, []byte) error
}

type Deleter interface {
	Delete(string) error
}

type Encoder interface {
	Encode() ([]byte, error)
	Length() int
}

// DBView defines the readonly view on database
// If readonly operation needed
// like - find and cursor
// then DatabseReader instance should be used
type DBView interface {
	Find(string, string) ([]byte, error)
	Cursor(string) <-chan []byte
	// All(string) ([][]byte, error)
	Len(string) uint64
}

// DB helps to define db layer
// it gives full access over store
// read and write both pervilege
type DB interface {
	io.Closer
	DBView

	Create(string, string, []byte) error
	CreateBulk(string, map[string][]byte) (int, error)
	Update(string, string, []byte) error
	Delete(string, string) error
}
