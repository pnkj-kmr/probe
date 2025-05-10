package probe

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type ScheduleEngine struct {
	engine  *actor.Engine
	process *safemap.SafeMap[int, *scheduleProcess]
}

func newScheduleEngine(e *actor.Engine, poll *PollEngine, opts ...OptFunc) (*ScheduleEngine, error) {
	schEngine := &ScheduleEngine{
		engine:  e,
		process: safemap.New[int, *scheduleProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.icmp {
		err := schEngine.multiSetup(ICMP, "icmp", poll, options)
		if err != nil {
			slog.Error("[SCHEDULE]", "err", err)
			return nil, err
		}
	}
	if options.snmp {
		err := schEngine.multiSetup(SNMP, "snmp", poll, options)
		if err != nil {
			slog.Error("[SCHEDULE]", "err", err)
			return nil, err
		}
	}
	return schEngine, nil
}

func (e *ScheduleEngine) Get(id int) (*scheduleProcess, bool) {
	return e.process.Get(id)
}

func (e *ScheduleEngine) Start() {
	e.process.ForEach(func(i int, s *scheduleProcess) {
		s.start()
	})
}

func (e *ScheduleEngine) Stop() {
	e.process.ForEach(func(i int, s *scheduleProcess) {
		s.stop()
	})
}

func (e *ScheduleEngine) setup(id int, name string, poll *PollEngine, interval time.Duration, options Opts) (err error) {
	slog.Info("setting up schedular", "name", name, "id", id)
	c, ok := poll.Get(id)
	if !ok {
		return ErrPOLLER
	}
	sch, err := newScheduleProcess(id, fmt.Sprintf("schedule/%s/%d", name, id), e.engine, interval, options.maxRestarts, c)
	if err != nil {
		return err
	}
	e.process.Set(id, sch)
	return nil
}

func (e *ScheduleEngine) multiSetup(id int, name string, poll *PollEngine, options Opts) (err error) {
	// setting up for 60 seconds
	err = e.setup(INTERVAL_60+id, name, poll, time.Second*time.Duration(INTERVAL_60), options)
	if err != nil {
		return
	}
	// setting up for 300 seconds
	err = e.setup(INTERVAL_300+id, name, poll, time.Second*time.Duration(INTERVAL_300), options)
	if err != nil {
		return
	}
	return
}
