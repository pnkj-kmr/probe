package boltdb

import (
	"context"
	"path/filepath"
	"probe/database"
	"time"

	bolt "go.etcd.io/bbolt"
)

var DBName string = "probe.db"

type store struct {
	ctx  context.Context
	db   *bolt.DB
	path string
}

func New(ctx context.Context, directory string) (database.Database, error) {
	if directory != "" {
		_, err := getOrCreateDir(directory)
		if err != nil {
			return nil, err
		}
	}
	path := filepath.Join(directory, DBName)
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}
	return &store{
		ctx:  ctx,
		db:   db,
		path: path,
	}, nil
}

// func (d *store) Path() string {
// 	return d.path
// }

func (d *store) Close() error {
	return d.db.Close()
}

func (d *store) Create(bucket, key string, data []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), data)
	})

}

func (d *store) Update(bucket, key string, data []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		return b.Put([]byte(key), data)
	})
}

func (d *store) Delete(bucket, key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		return b.Delete([]byte(key))
	})
}

func (d *store) Find(bucket, key string) (out []byte, err error) {
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

func (d *store) Cursor(bucket string) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		d.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				ch <- v
			}
			return nil
		})
	}()

	return ch
}

func (d *store) Len(bucket string) (total int) {
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		data := b.Stats()
		total = data.KeyN
		return nil
	})
	return
}

func (d *store) CreateBulk(bucket string, data map[string][]byte) (total int, err error) {
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

// func (d *store) all(bucket string) (out [][]byte, err error) {
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
