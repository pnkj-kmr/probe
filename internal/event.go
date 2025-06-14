package probe

import "github.com/anthdm/hollywood/actor"

type EventEngine struct {
	engine  *actor.Engine
	process *eventProcess
}

func newEventEngine(ctx *Context, options Opts) (err error) {
	eventEngine := &EventEngine{engine: ctx.Engine()}

	process, err := newEventProcess(0, "Event", ctx.Engine())
	if err != nil {
		return err
	}
	eventEngine.process = process

	ctx.WithEvent(eventEngine)
	return nil
}

// func (e *EventEngine) Get(id int) (any, bool) {
// 	return e.db.Get(id)
// }
