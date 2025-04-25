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

func newExportEngine(e *actor.Engine, opts ...OptFunc) (*ExportEngine, error) {
	engine := &ExportEngine{
		engine:  e,
		process: safemap.New[int, *exportProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	// assigning the export engine
	// TODO - need to create multi exporter as per poller count of other point
	if options.kafka {
		exprt, err := newExportProcess(options.context, KAFKA, "export/kafka", e, options.bucketSize, options.partition)
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
