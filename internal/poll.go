package probe

import (
	"fmt"
	"log/slog"
	M "probe/model"

	"probe/safemap"
)

type PollEngine struct {
	ctx     *Context
	process *safemap.SafeMap[string, *pollProcess]
}

func newPollEngine(ctx *Context, options Opts) (err error) {
	pollEngine := &PollEngine{
		ctx:     ctx,
		process: safemap.New[string, *pollProcess](),
	}

	if options.icmp {
		err := pollEngine.multiSetup(M.ICMP, options)
		if err != nil {
			slog.Error("[POLL] icmp", "err", err)
			return err
		}
	}
	if options.snmp {
		err := pollEngine.multiSetup(M.SNMP, options)
		if err != nil {
			slog.Error("[POLL] snmp", "err", err)
			return err
		}
	}
	ctx.WithPoll(pollEngine)
	return nil
}

func (e *PollEngine) Get(key string) (*pollProcess, bool) {
	return e.process.Get(key)
}

func (e *PollEngine) setup(name string, options Opts) (err error) {
	slog.Info("setting up poller", "name", name)
	db, export, event := e.ctx.DB(), e.ctx.Export(), e.ctx.Event()

	c, ok := db.GetDB(name)
	if !ok {
		return ErrDB
	}
	monitor, _ := event.Get(string(M.PROCESS))

	var exporter []M.Sender[any]
	if options.kafka {
		e, ok := export.Get(string(M.KAFKA))
		if !ok {
			return ErrNoKAFKA
		}
		exporter = append(exporter, e)
	}
	// checking db dump export if any
	ex, ok := export.Get(fmt.Sprintf("%s/%s_stat", M.DBDUMP, name))
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
	poller := newPollProcess(fmt.Sprintf("poll/%s", name), c, options.workers, exporter, finder, monitor)
	e.process.Set(name, poller)
	return nil
}

func (e *PollEngine) multiSetup(name M.ProbeType, options Opts) (err error) {
	// setting up for 60 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_60), options)
	if err != nil {
		return
	}
	// setting up for 300 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_300), options)
	if err != nil {
		return
	}
	return
}
