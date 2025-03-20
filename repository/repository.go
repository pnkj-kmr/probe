package repository

import (
	"context"
	"probe/repository/boltdb"
)

// Repository which returns the DB
type Repository struct {
	id        int
	bucket    string
	directory string
	dbname    string
	db        *boltdb.DB
}

func New(name string, id int) (*Repository, error) {
	db, err := boltdb.NewDB(context.Background(), "data", name)
	if err != nil {
		return nil, err
	}
	return &Repository{id: id, bucket: name, directory: "data", dbname: name, db: db}, nil
}

func (r *Repository) Consume() <-chan []byte {
	return r.db.Cursor(r.bucket)
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Create(key string, data []byte) error {
	return r.db.Create(r.bucket, key, data)
}

func (r *Repository) Find(key string) (out []byte, err error) {
	return r.db.Find(r.bucket, key)
}
