package probe

import (
	"fmt"
	"log/slog"
	M "probe/model"

	"probe/safemap"
)

type PollEngine struct {
	process *safemap.SafeMap[int, *pollProcess]
}

func newPollEngine(db *DBEngine, export *ExportEngine, opts ...OptFunc) (*PollEngine, error) {
	pollEngine := &PollEngine{
		process: safemap.New[int, *pollProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.icmp {
		err := pollEngine.multiSetup(M.ICMP, "icmp", db, export, options)
		if err != nil {
			slog.Error("[POLL]", "err", err)
			return nil, err
		}
	}
	if options.snmp {
		err := pollEngine.multiSetup(M.SNMP, "snmp", db, export, options)
		if err != nil {
			slog.Error("[POLL]", "err", err)
			return nil, err
		}
	}
	return pollEngine, nil
}

func (e *PollEngine) Get(id int) (*pollProcess, bool) {
	return e.process.Get(id)
}

func (e *PollEngine) setup(id int, name string, db *DBEngine, export *ExportEngine, options Opts) (err error) {
	slog.Info("setting up poller", "name", name, "id", id)
	c, ok := db.GetDB(id)
	if !ok {
		return ErrDB
	}
	var exporter *exportProcess
	if options.kafka {
		e, ok := export.Get(M.KAFKA)
		if !ok {
			return ErrNoKAFKA
		}
		exporter = e
	}
	var finder M.Finder = nil
	if options.snmp {
		f, ok := db.GetDB(M.CRED)
		if !ok {
			return ErrNoAuth
		}
		finder = f
	}
	poller := newPollProcess(id, fmt.Sprintf("poll/%s/%d", name, id), c, options.workers, exporter, finder)
	e.process.Set(id, poller)
	return nil
}

func (e *PollEngine) multiSetup(id int, name string, db *DBEngine, export *ExportEngine, options Opts) (err error) {
	// setting up for 60 seconds
	err = e.setup(M.INTERVAL_60+id, name, db, export, options)
	if err != nil {
		return
	}
	// setting up for 300 seconds
	err = e.setup(M.INTERVAL_300+id, name, db, export, options)
	if err != nil {
		return
	}
	return
}
