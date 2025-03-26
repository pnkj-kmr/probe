package probe

import (
	M "probe/model"
	"strconv"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type scheduleProcess struct {
	id          int
	name        string
	pid         *actor.PID
	repeater    actor.SendRepeater
	engine      *actor.Engine
	interval    time.Duration
	receiver    actor.Receiver
	maxRestarts int
}

func newScheduleProcess(id int, name string, e *actor.Engine, t time.Duration, maxRestarts int, receiver *pollProcess) (*scheduleProcess, error) {
	return &scheduleProcess{
		id: id, name: name, engine: e,
		interval:    t,
		receiver:    receiver,
		maxRestarts: maxRestarts,
	}, nil
}

func (p *scheduleProcess) start() {
	// log.Println("[SCHEDULE] starting...", p.name)
	p.pid = p.engine.SpawnFunc(
		p.receiver.Receive,
		p.name,
		actor.WithID(strconv.Itoa(p.id)),
		actor.WithMaxRestarts(p.maxRestarts),
	)
	p.repeater = p.engine.SendRepeat(
		p.pid,
		M.PollingBeat{Id: p.id, Name: p.name},
		p.interval,
	)
}

func (p *scheduleProcess) stop() error {
	// log.Println("[SCHEDULE] stopping...", p.name)
	p.repeater.Stop()
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	// log.Println("[SCHEDULE] stoppped", p.name)
	return nil
}
