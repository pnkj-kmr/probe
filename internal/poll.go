package probe

import (
	"github.com/anthdm/hollywood/safemap"
)

type PollEngine struct {
	process *safemap.SafeMap[int, *pollProcess]
}

func newPollEngine(db *DBEngine, export *ExportEngine, opts ...OptFunc) (*PollEngine, error) {
	engine := &PollEngine{
		process: safemap.New[int, *pollProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.icmp {
		c, ok := db.Get(ICMP)
		if !ok {
			return nil, ErrDB
		}
		export, ok := export.Get(KAFKA)
		if !ok {
			return nil, ErrNoKAFKA
		}
		engine.process.Set(ICMP, newPollProcess(ICMP, "icmp", c, options.workers, export))
	}
	if options.snmp {
		c, ok := db.Get(SNMP)
		if !ok {
			return nil, ErrDB
		}
		export, ok := export.Get(KAFKA)
		if !ok {
			return nil, ErrNoKAFKA
		}
		engine.process.Set(SNMP, newPollProcess(SNMP, "snmp", c, options.workers, export))
	}
	return engine, nil
}

func (e *PollEngine) Get(id int) (*pollProcess, bool) {
	return e.process.Get(id)
}
