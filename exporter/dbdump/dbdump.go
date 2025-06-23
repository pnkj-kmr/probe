package dbdump

import (
	"encoding/json"
	"fmt"
	"log/slog"
	M "probe/model"
)

type DBDump struct {
	name string
	db   M.DB
}

func NewDBDumpExporter(name string, db M.DB) *DBDump {
	return &DBDump{
		name: name,
		db:   db,
	}
}

func (dd *DBDump) Spin()          {}
func (dd *DBDump) Push() chan any { return nil }
func (dd *DBDump) Close() error   { return nil }

func (dd *DBDump) Export(data any) (err error) {
	switch d := data.(type) {
	case []any:
		slog.Info("dbdump received --->", "total", len(d))
		dump := make(map[string][]byte)
		for _, x := range d {
			rs, _ := json.Marshal(x) // TODO err case handling
			switch p := x.(type) {
			case M.ICMPRes:
				dump[p.Cid] = rs
			case M.SNMPRes:
				dump[p.Cid] = rs
			}
		}
		if len(dump) > 0 {
			dd.db.CreateBulk(dump)
		}
	default:
		fmt.Println("UNKNOWN DATA RECEIVED")
	}
	return
}
