package jsondb

import (
	"context"
	"path/filepath"
	"probe/database"

	simplejsondb "github.com/pnkj-kmr/simple-json-db"
)

var DBName string = "probe"

type store struct {
	ctx  context.Context
	db   *simplejsondb.DB
	path string
}

func New(ctx context.Context, directory string) (database.Database, error) {
	path := filepath.Join(directory, DBName)
	db, err := simplejsondb.New(path, nil)
	if err != nil {
		return nil, err
	}

	return &store{
		ctx:  ctx,
		db:   &db,
		path: path,
	}, nil
}

func (j *store) Close() error {
	return nil
}

func (j *store) Create(string, string, []byte) error {
	return nil
}
func (j *store) CreateBulk(string, map[string][]byte) (int, error) {
	return 0, nil
}

func (j *store) Update(string, string, []byte) error {
	return nil
}

func (j *store) Delete(string, string) error {
	return nil
}

func (j *store) Find(string, string) ([]byte, error) {
	return nil, nil
}

func (j *store) Cursor(string) <-chan []byte {
	return nil
}

func (j *store) Len(bucket string) (total uint64) {
	return
}
