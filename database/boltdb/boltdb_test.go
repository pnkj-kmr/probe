package boltdb_test

import (
	"context"
	"os"
	"path/filepath"
	"probe/database/boltdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func tempfile() string {
// 	f, err := os.CreateTemp("", "bolt-")
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := f.Close(); err != nil {
// 		panic(err)
// 	}
// 	if err := os.Remove(f.Name()); err != nil {
// 		panic(err)
// 	}
// 	return f.Name()
// }

func TestBolt(t *testing.T) {
	// path := tempfile()
	// defer os.RemoveAll(path)

	// db, err := boltdb.New(context.TODO(), path)
	// assert.Nil(t, err)
	// assert.NotNil(t, db)
	// db.Close()

	db, err := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	db.Close()

	db, err = boltdb.New(context.TODO(), "tmp")
	defer os.Remove(filepath.Join("tmp"))
	defer os.RemoveAll(filepath.Join("tmp", boltdb.DBName))
	assert.Nil(t, err)
	assert.NotNil(t, db)
	db.Close()
}

func TestBolt_Create(t *testing.T) {
	db, _ := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)

	val := []byte("v1")
	err := db.Create("b", "k1", val)
	assert.Nil(t, err)
	assert.Equal(t, val, []byte("v1"))
}

func TestBolt_Update(t *testing.T) {
	db, _ := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)

	val := []byte("v1")
	err := db.Update("b", "k1", val)
	assert.NotNil(t, err)
	err = db.Create("b", "k1", val)
	assert.Nil(t, err)
	assert.Equal(t, val, []byte("v1"))

	val = []byte("v2")
	err = db.Update("b", "k1", val)
	assert.Nil(t, err)
	assert.Equal(t, val, []byte("v2"))
}

func TestBolt_Delete(t *testing.T) {
	db, _ := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)

	val := []byte("v1")
	err := db.Delete("b", "k1")
	assert.Nil(t, err)

	err = db.Create("b", "k1", val)
	assert.Nil(t, err)
	assert.Equal(t, val, []byte("v1"))

	err = db.Delete("b", "k1")
	assert.Nil(t, err)

	data, _ := db.Find("b", "k1")
	assert.Nil(t, data)
}

func TestBolt_Find(t *testing.T) {
	db, _ := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)

	val := []byte("v1")
	_, err := db.Find("b", "k1")
	assert.NotNil(t, err)

	err = db.Create("b", "k1", val)
	assert.Nil(t, err)
	assert.Equal(t, val, []byte("v1"))

	data, err := db.Find("b", "k1")
	assert.Equal(t, data, val)

	err = db.Delete("b", "k1")
	assert.Nil(t, err)

	data, _ = db.Find("b", "k1")
	assert.Nil(t, data)
}

func TestBolt_Cursor(t *testing.T) {
	db, _ := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)

	val := []byte("v1")
	db.Create("b", "k1", val)
	db.Create("b", "k2", val)
	db.Create("b", "k3", val)

	data := db.Cursor("b")
	assert.NotNil(t, data)
	counter := 0
	for range data {
		counter += 1
	}
	assert.Equal(t, counter, 3)
}

func TestBolt_CreateBulk(t *testing.T) {
	db, _ := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)

	val := []byte("v1")
	data := map[string][]byte{
		"k1": val, "k2": val,
	}
	total, err := db.CreateBulk("b", data)
	assert.Equal(t, 2, total)
	assert.Nil(t, err)
}

func TestBolt_Len(t *testing.T) {
	db, _ := boltdb.New(context.TODO(), "")
	defer os.RemoveAll(boltdb.DBName)

	total := db.Len("b")
	assert.Equal(t, 0, total)

	val := []byte("v1")
	db.Create("b", "k1", val)
	db.Create("b", "k2", val)
	db.Create("b", "k3", val)

	total = db.Len("b")
	assert.Equal(t, 3, total)
}
