package probe

import (
	M "probe/model"
	"probe/safemap"
)

type EventEngine struct {
	process *safemap.SafeMap[string, *eventProcess]
}

func newEventEngine(ctx *Context, _ Opts) (err error) {
	eventEngine := &EventEngine{
		process: safemap.New[string, *eventProcess](),
	}

	// process, err := newEventProcess(string(M.EVENT), ctx.Engine())
	// if err != nil {
	// 	return err
	// }
	// eventEngine.process.Set(string(M.EVENT), process)

	p, err := newEventProcess(string(M.PROCESS), ctx.Engine(), ctx.DB())
	if err != nil {
		return err
	}
	eventEngine.process.Set(string(M.PROCESS), p)

	ctx.WithEvent(eventEngine)
	return nil
}

func (e *EventEngine) Get(k string) (*eventProcess, bool) {
	return e.process.Get(k)
}

func (e *EventEngine) Start() {
	e.process.ForEach(func(i string, e *eventProcess) {
		e.Start()
	})
}

func (e *EventEngine) Stop() {
	e.process.ForEach(func(i string, e *eventProcess) {
		e.Stop()
	})
}
