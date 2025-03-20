package probe

import (
	"fmt"
	M "probe/model"
	"strconv"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type ScheduleEngine struct {
	engine    *actor.Engine
	scheduler *safemap.SafeMap[int, *schedule]
}

func NewScheduleEngine(poll *PollEngine, opts ...OptFunc) (*ScheduleEngine, error) {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return nil, err
	}
	engine := &ScheduleEngine{
		engine:    e,
		scheduler: safemap.New[int, *schedule](),
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
		sch, err := newSchedule(ICMP, "icmp", e, options.interval, options.maxRestarts, c)
		if err != nil {
			return nil, err
		}
		engine.scheduler.Set(ICMP, sch)
	}
	if options.snmp {
		c, ok := poll.Get(SNMP)
		if !ok {
			return nil, ErrPOLLER
		}
		sch, err := newSchedule(SNMP, "snmp", e, options.interval, options.maxRestarts, c)
		if err != nil {
			return nil, err
		}
		engine.scheduler.Set(SNMP, sch)
	}
	return engine, nil
}

func (e *ScheduleEngine) Get(id int) (*schedule, bool) {
	return e.scheduler.Get(id)
}

type schedule struct {
	id          int
	name        string
	pid         *actor.PID
	repeater    actor.SendRepeater
	engine      *actor.Engine
	interval    time.Duration
	receiver    actor.Receiver
	maxRestarts int
}

func newSchedule(id int, name string, e *actor.Engine, t time.Duration, maxRestarts int, receiver *process) (*schedule, error) {
	return &schedule{id: id, name: name, engine: e, interval: t, receiver: receiver, maxRestarts: maxRestarts}, nil
}

func (s *schedule) Start() {
	fmt.Println("schedule---starting------", s.name)
	s.pid = s.engine.SpawnFunc(s.receiver.Receive, s.name, actor.WithID(strconv.Itoa(s.id)), actor.WithMaxRestarts(s.maxRestarts))
	s.repeater = s.engine.SendRepeat(s.pid, M.PollingBeat{Id: s.id, Name: s.name}, s.interval)
}

func (s *schedule) Stop() error {
	fmt.Println("schedule---stopping------", s.name)
	s.repeater.Stop()
	ctx := s.engine.Poison(s.pid)
	<-ctx.Done()
	fmt.Println("schedule---stopped------", s.name)
	return nil
}
