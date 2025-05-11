package probe

import (
	"fmt"
	"log/slog"
	M "probe/model"

	"probe/safemap"

	"github.com/anthdm/hollywood/actor"
)

type ExportEngine struct {
	// workers  int
	engine  *actor.Engine
	process *safemap.SafeMap[int, *exportProcess]
}

func newExportEngine(e *actor.Engine, opts ...OptFunc) (*ExportEngine, error) {
	exportEngine := &ExportEngine{
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
		err := exportEngine.setup(M.KAFKA, "kafka", options)
		if err != nil {
			slog.Error("[EXPORT]", "err", err)
			return nil, err
		}
	}
	// add a webhook exportProcess if needed

	return exportEngine, nil
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

func (e *ExportEngine) setup(id int, name string, options Opts) (err error) {
	slog.Info("setting up exporter", "name", name, "id", id)
	exprt, err := newExportProcess(options.context, id, fmt.Sprintf("export/%s/%d", name, id), e.engine, options.maxBucketSize, options.totalPartitions)
	if err != nil {
		return err
	}
	e.process.Set(id, exprt)
	return nil
}
