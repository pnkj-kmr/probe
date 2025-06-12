package repository

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"probe/repository/boltdb"
)

const DEFAULT_DIRECTORY string = "data"

// Repository which returns the DB
type Repository struct {
	//table name into a repo
	bucket string
	// path where repo located cwd + directory
	directory string
	// name of db file
	dbname string
	// instance of db lib
	db *boltdb.DB
}

func New(name, directory string) (*Repository, error) {
	pathList := strings.Split(name, "/")
	if len(pathList) > 1 {
		name = pathList[len(pathList)-1]
		directory = filepath.Join(pathList[:len(pathList)-1]...)
	}
	if directory != "" {
		directory = filepath.Join(DEFAULT_DIRECTORY, directory)
	} else {
		directory = DEFAULT_DIRECTORY
	}
	// slog.Info("creating new db", "directory", directory, "name", name)
	db, err := boltdb.NewDB(context.Background(), directory, fmt.Sprintf("%s.db", name))
	if err != nil {
		return nil, err
	}
	return &Repository{bucket: name, directory: directory, dbname: fmt.Sprintf("%s.db", name), db: db}, nil
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

func (r *Repository) DeleteAll() error {
	return r.db.DeleteAll(r.bucket)
}

func (r *Repository) Find(key string) (out []byte, err error) {
	return r.db.Find(r.bucket, key)
}

func (r *Repository) CreateBulk(data map[string][]byte) (total int, err error) {
	return r.db.CreateBulk(r.bucket, data)
}
