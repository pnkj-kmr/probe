package probe

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type ScheduleEngine struct {
	engine  *actor.Engine
	process *safemap.SafeMap[int, *scheduleProcess]
}

func newScheduleEngine(e *actor.Engine, poll *PollEngine, opts ...OptFunc) (*ScheduleEngine, error) {
	engine := &ScheduleEngine{
		engine:  e,
		process: safemap.New[int, *scheduleProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.icmp {
		c, ok := poll.Get(ICMP)
		if !ok {
			return nil, ErrPOLLER
		}
		sch, err := newScheduleProcess(ICMP, "schedule/icmp", e, options.interval, options.maxRestarts, c)
		if err != nil {
			return nil, err
		}
		engine.process.Set(ICMP, sch)
	}
	if options.snmp {
		c, ok := poll.Get(SNMP)
		if !ok {
			return nil, ErrPOLLER
		}
		sch, err := newScheduleProcess(SNMP, "schedule/snmp", e, options.interval, options.maxRestarts, c)
		if err != nil {
			return nil, err
		}
		engine.process.Set(SNMP, sch)
	}
	return engine, nil
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
