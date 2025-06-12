package probe

import (
	M "probe/model"
	"probe/safemap"

	"github.com/anthdm/hollywood/actor"
)

type APIEngine struct {
	engine  *actor.Engine
	process *safemap.SafeMap[M.ProbeType, *apiProcess]
}

func newAPIEngine(e *actor.Engine, db *DBEngine, opts ...OptFunc) (*APIEngine, error) {
	apiEngine := &APIEngine{
		engine:  e,
		process: safemap.New[M.ProbeType, *apiProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.api {
		p, err := newApiProcess(M.API, e, db)
		if err != nil {
			return nil, err
		}
		apiEngine.process.Set(M.API, p)
	}

	return apiEngine, nil
}

func (e *APIEngine) Get(k M.ProbeType) (*apiProcess, bool) {
	return e.process.Get(k)
}

func (e *APIEngine) Start() {
	e.process.ForEach(func(i M.ProbeType, s *apiProcess) {
		s.start()
	})
}

func (e *APIEngine) Stop() {
	e.process.ForEach(func(i M.ProbeType, s *apiProcess) {
		s.stop()
	})
}
