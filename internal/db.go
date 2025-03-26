package probe

import (
	"probe/repository"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type DBEngine struct {
	engine  *actor.Engine
	db      *safemap.SafeMap[int, *repository.Repository]
	process *safemap.SafeMap[int, *dbProcess]
}

func newDBEngine(opts ...OptFunc) (*DBEngine, error) {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return nil, err
	}

	engine := &DBEngine{
		engine:  e,
		db:      safemap.New[int, *repository.Repository](),
		process: safemap.New[int, *dbProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.icmp {
		db, err := repository.New("icmp", ICMP)
		if err != nil {
			return nil, err
		}
		engine.db.Set(ICMP, db)

		p, err := newDBProcess(ICMP, "icmp", e, db)
		if err != nil {
			return nil, err
		}
		engine.process.Set(ICMP, p)
	}
	if options.snmp {
		db, err := repository.New("snmp", SNMP)
		if err != nil {
			return nil, err
		}
		engine.db.Set(SNMP, db)

		p, err := newDBProcess(SNMP, "snmp", e, db)
		if err != nil {
			return nil, err
		}
		engine.process.Set(SNMP, p)
	}
	return engine, nil
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
