package probe

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type APIEngine struct {
	engine  *actor.Engine
	process *safemap.SafeMap[int, *apiProcess]
}

func newAPIEngine(opts ...OptFunc) (*APIEngine, error) {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return nil, err
	}

	engine := &APIEngine{
		engine:  e,
		process: safemap.New[int, *apiProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.api {
		p, err := newApiProcess(API, "api", e)
		if err != nil {
			return nil, err
		}
		engine.process.Set(API, p)
	}

	return engine, nil
}

func (e *APIEngine) Get(id int) (*apiProcess, bool) {
	return e.process.Get(id)
}

func (e *APIEngine) Start() {
	e.process.ForEach(func(i int, s *apiProcess) {
		s.start()
	})
}

func (e *APIEngine) Stop() {
	e.process.ForEach(func(i int, s *apiProcess) {
		s.stop()
	})
}
