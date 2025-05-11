package probe

import (
	"fmt"
	"log/slog"
	"path/filepath"
	M "probe/model"
	"probe/repository"
	"probe/safemap"
	"strconv"

	"github.com/anthdm/hollywood/actor"
)

type DBEngine struct {
	engine  *actor.Engine
	db      *safemap.SafeMap[int, *repository.Repository]
	process *safemap.SafeMap[int, *dbProcess]
}

func newDBEngine(e *actor.Engine, opts ...OptFunc) (*DBEngine, error) {
	dbEngine := &DBEngine{
		engine:  e,
		db:      safemap.New[int, *repository.Repository](),
		process: safemap.New[int, *dbProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	// initialising db
	var err error

	err = dbEngine.setup(M.AUTH_PROFILE, "auth", "")
	if err != nil {
		return nil, err
	}
	err = dbEngine.setup(M.CONFIG, "config", "")
	if err != nil {
		return nil, err
	}
	err = dbEngine.setup(M.ENV, "env", "")
	if err != nil {
		return nil, err
	}
	err = dbEngine.multiSetup(M.ICMP, "icmp")
	if err != nil {
		return nil, err
	}
	err = dbEngine.multiSetup(M.SNMP, "snmp")
	if err != nil {
		return nil, err
	}

	return dbEngine, nil
}

func (e *DBEngine) GetDB(id int) (*repository.Repository, bool) {
	return e.db.Get(id)
}

func (e *DBEngine) GetProcess(id int) (*dbProcess, bool) {
	return e.process.Get(id)
}

func (e *DBEngine) Start() {
	e.process.ForEach(func(i int, s *dbProcess) {
		s.start()
	})
}

func (e *DBEngine) Stop() {
	e.process.ForEach(func(i int, s *dbProcess) {
		s.stop()
	})
}

func (e *DBEngine) setup(id int, name, directory string) (err error) {
	slog.Info("setting up db", "name", name, "directory", directory, "id", id)
	db, err := repository.New(id, name, directory)
	if err != nil {
		slog.Error("[DB] repo error", "err", err)
		return err
	}
	e.db.Set(id, db)

	p, err := newDBProcess(id, fmt.Sprintf("db/%s/%d", filepath.Join(directory, name), id), e.engine, db)
	if err != nil {
		slog.Error("[DB] process error", "err", err)
		return err
	}
	e.process.Set(id, p)
	return
}

func (e *DBEngine) multiSetup(id int, name string) (err error) {
	// setting up for 60 seconds
	err = e.setup(M.INTERVAL_60+id, strconv.Itoa(M.INTERVAL_60), name)
	if err != nil {
		return
	}
	// setting up for 300 seconds
	err = e.setup(M.INTERVAL_300+id, strconv.Itoa(M.INTERVAL_300), name)
	if err != nil {
		return
	}
	return
}
