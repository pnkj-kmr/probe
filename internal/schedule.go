package probe

import (
	"fmt"
	"log/slog"
	M "probe/model"
	"time"

	"probe/safemap"

	"github.com/anthdm/hollywood/actor"
)

type ScheduleEngine struct {
	engine  *actor.Engine
	process *safemap.SafeMap[M.ProbeType, *scheduleProcess]
}

func newScheduleEngine(ctx *Context, options Opts) (err error) {
	schEngine := &ScheduleEngine{
		engine:  ctx.Engine(),
		process: safemap.New[M.ProbeType, *scheduleProcess](),
	}

	if options.icmp {
		err := schEngine.multiSetup(M.ICMP, ctx.Poll(), options)
		if err != nil {
			slog.Error("[SCHEDULE]", "err", err)
			return err
		}
	}
	if options.snmp {
		err := schEngine.multiSetup(M.SNMP, ctx.Poll(), options)
		if err != nil {
			slog.Error("[SCHEDULE]", "err", err)
			return err
		}
	}

	ctx.WithSchdule(schEngine)
	return nil
}

func (e *ScheduleEngine) Get(name M.ProbeType) (*scheduleProcess, bool) {
	return e.process.Get(name)
}

func (e *ScheduleEngine) Start() {
	e.process.ForEach(func(i M.ProbeType, s *scheduleProcess) {
		s.start()
	})
}

func (e *ScheduleEngine) Stop() {
	e.process.ForEach(func(i M.ProbeType, s *scheduleProcess) {
		s.stop()
	})
}

func (e *ScheduleEngine) setup(name string, poll *PollEngine, interval time.Duration, options Opts) (err error) {
	slog.Info("setting up schedular", "name", name)
	c, ok := poll.Get(M.ProbeType(name))
	if !ok {
		return ErrPOLLER
	}
	sch, err := newScheduleProcess(fmt.Sprintf("schedule/%s", name), e.engine, interval, options.maxRestarts, c)
	if err != nil {
		return err
	}
	e.process.Set(M.ProbeType(name), sch)
	return nil
}

func (e *ScheduleEngine) multiSetup(name M.ProbeType, poll *PollEngine, options Opts) (err error) {
	// setting up for 60 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_60), poll, time.Second*time.Duration(M.INTERVAL_60), options)
	if err != nil {
		return
	}
	// setting up for 300 seconds
	err = e.setup(fmt.Sprintf("%s/%d", name, M.INTERVAL_300), poll, time.Second*time.Duration(M.INTERVAL_300), options)
	if err != nil {
		return
	}
	return
}
