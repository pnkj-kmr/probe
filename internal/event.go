package probe

import "github.com/anthdm/hollywood/actor"

type EventEngine struct {
	engine  *actor.Engine
	process *eventProcess
}

func newEventEngine(e *actor.Engine, opts ...OptFunc) (*EventEngine, error) {

	engine := &EventEngine{engine: e}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	process, err := newEventProcess(0, "Event", e)
	if err != nil {
		return nil, err
	}
	engine.process = process

	return engine, nil
}

// func (e *EventEngine) Get(id int) (any, bool) {
// 	return e.db.Get(id)
// }
