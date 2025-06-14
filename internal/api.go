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

func newAPIEngine(ctx *Context, options Opts) (err error) {
	apiEngine := &APIEngine{
		engine:  ctx.Engine(),
		process: safemap.New[M.ProbeType, *apiProcess](),
	}

	if options.api {
		p, err := newApiProcess(M.API, ctx.Engine(), ctx.DB())
		if err != nil {
			return err
		}
		apiEngine.process.Set(M.API, p)
	}

	ctx.WithAPI(apiEngine)
	return nil
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
