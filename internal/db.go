package probe

import (
	"probe/repository"

	"github.com/anthdm/hollywood/safemap"
)

type DBEngine struct {
	db *safemap.SafeMap[int, *repository.Repository]
}

func NewDBEngine(opts ...OptFunc) (*DBEngine, error) {
	engine := &DBEngine{
		db: safemap.New[int, *repository.Repository](),
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
	}
	if options.snmp {
		db, err := repository.New("snmp", SNMP)
		if err != nil {
			return nil, err
		}
		engine.db.Set(SNMP, db)
	}
	return engine, nil
}

func (e *DBEngine) Get(id int) (*repository.Repository, bool) {
	return e.db.Get(id)
}
