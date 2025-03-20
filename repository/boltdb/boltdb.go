package boltdb

import (
	"context"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

var DBName string = "probe.db"

type DB struct {
	ctx  context.Context
	db   *bolt.DB
	path string
}

func New(ctx context.Context, directory string) (*DB, error) {
	return NewDB(ctx, directory, DBName)
}

func NewDB(ctx context.Context, directory, db_name string) (*DB, error) {
	if directory != "" {
		_, err := getOrCreateDir(directory)
		if err != nil {
			return nil, err
		}
	}
	if db_name == "" {
		db_name = DBName
	}
	path := filepath.Join(directory, db_name)
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}
	return &DB{
		ctx:  ctx,
		db:   db,
		path: path,
	}, nil
}

// func (d *DB) Path() string {
// 	return d.path
// }

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Create(bucket, key string, data []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), data)
	})

}

func (d *DB) Update(bucket, key string, data []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		return b.Put([]byte(key), data)
	})
}

func (d *DB) Delete(bucket, key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		return b.Delete([]byte(key))
	})
}

func (d *DB) Find(bucket, key string) (out []byte, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		v := b.Get([]byte(key))
		out = v
		return nil
	})
	return
}

func (d *DB) Cursor(bucket string) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		d.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b != nil {
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					ch <- v
				}
			}
			return nil
		})
	}()

	return ch
}

func (d *DB) Len(bucket string) (total uint64) {
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		data := b.Stats()
		total = uint64(data.KeyN)
		return nil
	})
	return
}

func (d *DB) CreateBulk(bucket string, data map[string][]byte) (total int, err error) {
	err = d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		var _err error
		for k, v := range data {
			e := b.Put([]byte(k), v)
			if e != nil {
				_err = e
			} else {
				total += 1
			}
		}
		return _err
	})
	return
}

// func (d *DB) all(bucket string) (out [][]byte, err error) {
// 	// var out [][]byte
// 	err = d.db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte(bucket)) // Assume bucket exists and has keys
// 		c := b.Cursor()
// 		for k, v := c.First(); k != nil; k, v = c.Next() {
// 			// fmt.Printf("key=%s, value=%s\n", k, v)
// 			out = append(out, v)
// 		}
// 		return nil
// 	})
// 	return
// }
