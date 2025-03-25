package probe

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type ExportEngine struct {
	// workers  int
	engine  *actor.Engine
	process *safemap.SafeMap[int, *exportProcess]
}

func newExportEngine(workers int, opts ...OptFunc) (*ExportEngine, error) {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return nil, err
	}
	engine := &ExportEngine{
		engine:  e,
		process: safemap.New[int, *exportProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	if options.icmp || options.snmp {
		exprt, err := newExportProcess(KAFKA, "kafka", e, workers)
		if err != nil {
			return nil, err
		}
		engine.process.Set(KAFKA, exprt)
	}
	// add a webhook exportProcess if needed

	return engine, nil
}

func (e *ExportEngine) Get(id int) (*exportProcess, bool) {
	return e.process.Get(id)
}

func (e *ExportEngine) Start() {
	e.process.ForEach(func(i int, e *exportProcess) {
		e.Start()
	})
}

func (e *ExportEngine) Stop() {
	e.process.ForEach(func(i int, e *exportProcess) {
		e.Stop()
	})
}
