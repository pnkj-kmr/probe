package probe

import (
	"fmt"
	"log/slog"
	M "probe/model"

	"probe/safemap"
)

type PollEngine struct {
	process *safemap.SafeMap[M.ProbeType, *pollProcess]
}

func newPollEngine(db *DBEngine, export *ExportEngine, opts ...OptFunc) (*PollEngine, error) {
	pollEngine := &PollEngine{
		process: safemap.New[M.ProbeType, *pollProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.icmp {
		err := pollEngine.multiSetup(M.ICMP, db, export, options)
		if err != nil {
			slog.Error("[POLL] icmp", "err", err)
			return nil, err
		}
	}
	if options.snmp {
		err := pollEngine.multiSetup(M.SNMP, db, export, options)
		if err != nil {
			slog.Error("[POLL] snmp", "err", err)
			return nil, err
		}
	}
	return pollEngine, nil
}

func (e *PollEngine) Get(id M.ProbeType) (*pollProcess, bool) {
	return e.process.Get(id)
}

func (e *PollEngine) setup(name string, db *DBEngine, export *ExportEngine, options Opts) (err error) {
	slog.Info("setting up poller", "name", name)
	c, ok := db.GetDB(name)
	if !ok {
		return ErrDB
	}
	var exporter []M.Sender[any]
	if options.kafka {
		e, ok := export.Get(string(M.KAFKA))
		if !ok {
			return ErrNoKAFKA
		}
		exporter = append(exporter, e)
	}
	// checking db dump export if any
	ex, ok := export.Get(name)
	if ok {
		exporter = append(exporter, ex)
	}

	var finder M.Finder = nil
	if options.snmp {
		f, ok := db.GetDB(string(M.CRED))
		if !ok {
			return ErrNoAuth
		}
		finder = f
	}
	poller := newPollProcess(fmt.Sprintf("poll/%s", name), c, options.workers, exporter, finder)
	e.process.Set(M.ProbeType(name), poller)
	return nil
}

func (e *PollEngine) multiSetup(name M.ProbeType, db *DBEngine, export *ExportEngine, options Opts) (err error) {
	// setting up for 60 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_60), db, export, options)
	if err != nil {
		return
	}
	// setting up for 300 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_300), db, export, options)
	if err != nil {
		return
	}
	return
}
