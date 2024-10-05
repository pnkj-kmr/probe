package jsondb

// import (
// 	"context"
// 	"probe/database"

// 	simplejsondb "github.com/pnkj-kmr/simple-json-db"
// )

// type jdb struct {
// 	ctx context.Context
// 	db  *simplejsondb.DB
// }

// func New(ctx context.Context) database.Database {
// 	// TODO extra infomation for DB
// 	db, _ := simplejsondb.New("data/probe", nil)

// 	return &jdb{
// 		ctx: ctx,
// 		db:  &db,
// 	}
// }

// func (j *jdb) Close() error {
// 	return nil
// }

// func (j *jdb) All(bucket string) (out [][]byte, err error) {
// 	return
// }
// func (j *jdb) Find(bucket, key string) ([]byte, error) {
// 	return nil, nil
// }
// func (j *jdb) Create(bucket string, data [][]byte) (int, error) {
// 	return 0, nil
// }
// func (j *jdb) Update(bucket, key string, data []byte) ([]byte, error) {
// 	return nil, nil
// }
// func (j *jdb) Delete(bucket, key string) ([]byte, error) {
// 	return nil, nil
// }

// func (j *jdb) Cursor(string) (<-chan []byte, error) {
// 	return nil, nil
// }
