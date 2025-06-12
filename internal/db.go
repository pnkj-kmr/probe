package probe

import (
	"fmt"
	"log/slog"
	M "probe/model"
	"probe/repository"
	"probe/safemap"

	"github.com/anthdm/hollywood/actor"
)

type DBEngine struct {
	engine *actor.Engine
	db     *safemap.SafeMap[string, *repository.Repository]
	// process *safemap.SafeMap[int, *dbProcess]
}

func newDBEngine(e *actor.Engine, opts ...OptFunc) (*DBEngine, error) {
	dbEngine := &DBEngine{
		engine: e,
		db:     safemap.New[string, *repository.Repository](),
		// process: safemap.New[int, *dbProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	var err error

	err = dbEngine.setup(string(M.CRED), "")
	if err != nil {
		return nil, err
	}
	err = dbEngine.setup(string(M.CONFIG), "")
	if err != nil {
		return nil, err
	}
	err = dbEngine.setup(string(M.ENV), "")
	if err != nil {
		return nil, err
	}
	err = dbEngine.multiSetup(string(M.ICMP))
	if err != nil {
		return nil, err
	}
	err = dbEngine.multiSetup(string(M.SNMP))
	if err != nil {
		return nil, err
	}

	return dbEngine, nil
}

func (e *DBEngine) GetDB(name string) (*repository.Repository, bool) {
	return e.db.Get(name)
}

// func (e *DBEngine) GetProcess(id int) (*dbProcess, bool) {
// 	return e.process.Get(id)
// }

// func (e *DBEngine) Start() {
// 	e.process.ForEach(func(i int, s *dbProcess) {
// 		s.start()
// 	})
// }

// func (e *DBEngine) Stop() {
// 	e.process.ForEach(func(i int, s *dbProcess) {
// 		s.stop()
// 	})
// }

func (e *DBEngine) setup(name, directory string) (err error) {
	slog.Info("setting...", "name", name, "directory", directory)
	db, err := repository.New(name, directory)
	if err != nil {
		slog.Error("[DB] repo error", "err", err)
		return err
	}
	e.db.Set(name, db)

	// p, err := newDBProcess(id, fmt.Sprintf("db/%s/%d", filepath.Join(directory, name), id), e.engine, db)
	// if err != nil {
	// 	slog.Error("[DB] process error", "err", err)
	// 	return err
	// }
	// e.process.Set(id, p)
	return
}

func (e *DBEngine) multiSetup(name string) (err error) {
	// setting up for 60 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_60), "")
	if err != nil {
		return
	}
	err = e.setup(fmt.Sprintf("%s/%d_stat", name, M.INTERVAL_60), "")
	if err != nil {
		return
	}

	// setting up for 300 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_300), "")
	if err != nil {
		return
	}
	err = e.setup(fmt.Sprintf("%s/%d_stat", name, M.INTERVAL_300), "")
	if err != nil {
		return
	}
	return
}
