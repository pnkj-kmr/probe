package probe

import (
	"fmt"
	"log/slog"
	M "probe/model"
	"probe/repository"

	"probe/safemap"

	"github.com/anthdm/hollywood/actor"
)

type ExportEngine struct {
	// workers  int
	db      *DBEngine
	engine  *actor.Engine
	process *safemap.SafeMap[string, *exportProcess]
}

func newExportEngine(e *actor.Engine, db *DBEngine, opts ...OptFunc) (*ExportEngine, error) {
	exportEngine := &ExportEngine{
		engine: e, db: db,
		process: safemap.New[string, *exportProcess](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	// assigning the export engine
	// TODO - need to create multi exporter as per poller count of other point
	if options.kafka {
		err := exportEngine.setup(string(M.KAFKA), "", options)
		if err != nil {
			slog.Error("[EXPORT]", "err", err)
			return nil, err
		}
	}
	if options.icmp {
		err := exportEngine.setup(string(M.DBDUMP), fmt.Sprintf("%s/%d_stat", M.ICMP, M.INTERVAL_60), options)
		if err != nil {
			slog.Error("[EXPORT] icmp stat", "err", err)
			return nil, err
		}
		err = exportEngine.setup(string(M.DBDUMP), fmt.Sprintf("%s/%d_stat", M.ICMP, M.INTERVAL_300), options)
		if err != nil {
			slog.Error("[EXPORT] icmp stat", "err", err)
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

func (e *ExportEngine) setup(name, purpose string, options Opts) (err error) {
	slog.Info("setting up exporter", "name", name, "purpose", purpose)
	var exporterName string
	if purpose != "" {
		exporterName = fmt.Sprintf("%s/%s", name, purpose)
	} else {
		exporterName = name
	}

	var db *repository.Repository
	if purpose != "" {
		d, ok := e.db.GetDB(purpose)
		if !ok {
			slog.Error("db not found", "purpose", purpose, "db", d)
			return ErrDB
		}
		db = d
	}

	exprt, err := newExportProcess(
		options.context, exporterName, e.engine, db,
		options.maxBucketSize, options.totalPartitions,
	)
	if err != nil {
		return err
	}

	e.process.Set(exporterName, exprt)
	return nil
}
