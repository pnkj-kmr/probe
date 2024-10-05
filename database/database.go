package database

import "io"

// DatabaseReader defines the readonly view on database
// If readonly operation needed
// like - find and cursor
// then DatabseReader instance should be used
type DatabaseReader interface {
	Find(string, string) ([]byte, error)
	Cursor(string) <-chan []byte
	// All(string) ([][]byte, error)
	Len(string) uint64
}

// Database helps to define db layer
// it gives full access over database
// read and write both pervilege
type Database interface {
	io.Closer
	DatabaseReader

	Create(string, string, []byte) error
	CreateBulk(string, map[string][]byte) (int, error)
	Update(string, string, []byte) error
	Delete(string, string) error
}
