package repository

import (
	"context"
	"fmt"

	"probe/repository/boltdb"
)

// Repository which returns the DB
type Repository struct {
	// assign the id to repo instance
	id int
	//table name into a repo
	bucket string
	// path where repo located cwd + directory
	directory string
	// name of db file
	dbname string
	// instance of db lib
	db *boltdb.DB
}

func New(name string, id int) (*Repository, error) {
	db, err := boltdb.NewDB(context.Background(), "data", fmt.Sprintf("%s.db", name))
	if err != nil {
		return nil, err
	}
	return &Repository{id: id, bucket: name, directory: "data", dbname: fmt.Sprintf("%s.db", name), db: db}, nil
}

func (r *Repository) Receive() <-chan []byte {
	return r.db.Cursor(r.bucket)
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Count() uint64 {
	return r.db.Len(r.bucket)
}

func (r *Repository) Create(key string, data []byte) error {
	return r.db.Create(r.bucket, key, data)
}

func (r *Repository) Update(key string, data []byte) error {
	return r.db.Update(r.bucket, key, data)
}

func (r *Repository) Delete(key string) error {
	return r.db.Delete(r.bucket, key)
}

func (r *Repository) Find(key string) (out []byte, err error) {
	return r.db.Find(r.bucket, key)
}
