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
}

func newDBEngine(ctx *Context, _ Opts) (err error) {
	dbEngine := &DBEngine{
		engine: ctx.Engine(),
		db:     safemap.New[string, *repository.Repository](),
	}

	err = dbEngine.setup(string(M.CRED), "")
	if err != nil {
		return err
	}
	err = dbEngine.setup(string(M.CONFIG), "")
	if err != nil {
		return err
	}
	err = dbEngine.setup(string(M.ENV), "")
	if err != nil {
		return err
	}
	err = dbEngine.setup(string(M.PROCESS), "")
	if err != nil {
		return err
	}
	err = dbEngine.multiSetup(string(M.ICMP))
	if err != nil {
		return err
	}
	err = dbEngine.multiSetup(string(M.SNMP))
	if err != nil {
		return err
	}

	ctx.WithDB(dbEngine)
	return
}

func (e *DBEngine) GetDB(name string) (*repository.Repository, bool) {
	return e.db.Get(name)
}

func (e *DBEngine) setup(name, directory string) (err error) {
	slog.Info("setting...", "name", name, "directory", directory)
	db, err := repository.New(name, directory)
	if err != nil {
		slog.Error("[DB] repo error", "err", err)
		return err
	}
	e.db.Set(name, db)
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
