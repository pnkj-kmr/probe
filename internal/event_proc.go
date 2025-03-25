package probe

import (
	"github.com/anthdm/hollywood/actor"
)

type eventProcess struct {
	id     int
	name   string
	engine *actor.Engine
	// pid      *actor.PID
}

func newEventProcess(id int, name string, e *actor.Engine) (*eventProcess, error) {
	return &eventProcess{
		id: id, name: name, engine: e,
	}, nil
}
