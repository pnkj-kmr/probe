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
	process *safemap.SafeMap[string, *exportProcess]
}

func newExportEngine(e *actor.Engine, opts ...OptFunc) (*ExportEngine, error) {
	exportEngine := &ExportEngine{
		engine:  e,
		process: safemap.New[string, *exportProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	// assigning the export engine
	// TODO - need to create multi exporter as per poller count of other point
	if options.kafka {
		err := exportEngine.setup(string(M.KAFKA), options)
		if err != nil {
			slog.Error("[EXPORT]", "err", err)
			return nil, err
		}
	}
	// add a webhook exportProcess if needed

	return exportEngine, nil
}

func (e *ExportEngine) Get(k string) (*exportProcess, bool) {
	return e.process.Get(k)
}

func (e *ExportEngine) Start() {
	e.process.ForEach(func(i string, e *exportProcess) {
		e.Start()
	})
}

func (e *ExportEngine) Stop() {
	e.process.ForEach(func(i string, e *exportProcess) {
		e.Stop()
	})
}

func (e *ExportEngine) setup(name string, options Opts) (err error) {
	slog.Info("setting up exporter", "name", name)
	exprt, err := newExportProcess(options.context, fmt.Sprintf("export/%s", name), e.engine, options.maxBucketSize, options.totalPartitions)
	if err != nil {
		return err
	}
	e.process.Set(name, exprt)
	return nil
}
