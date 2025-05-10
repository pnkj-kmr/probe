package probe

import "github.com/anthdm/hollywood/actor"

type EventEngine struct {
	engine  *actor.Engine
	process *eventProcess
}

func newEventEngine(e *actor.Engine, opts ...OptFunc) (*EventEngine, error) {

	eventEngine := &EventEngine{engine: e}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	process, err := newEventProcess(0, "Event", e)
	if err != nil {
		return nil, err
	}
	eventEngine.process = process

	return eventEngine, nil
}

// func (e *EventEngine) Get(id int) (any, bool) {
// 	return e.db.Get(id)
// }
