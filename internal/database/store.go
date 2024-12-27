package database

import "io"

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
