package probe

import "github.com/anthdm/hollywood/actor"

type dbProcess struct {
	id     int
	name   string
	engine *actor.Engine
	// pid      *actor.PID
}

func newDBProcess(id int, name string, e *actor.Engine) (*dbProcess, error) {
	return &dbProcess{
		id: id, name: name, engine: e,
	}, nil
}
