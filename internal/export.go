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
	ctx     *Context
	db      *DBEngine
	engine  *actor.Engine
	process *safemap.SafeMap[string, *exportProcess]
}

func newExportEngine(ctx *Context, opts Opts) (err error) {
	exportEngine := &ExportEngine{
		ctx:    ctx,
		engine: ctx.Engine(), db: ctx.DB(),
		process: safemap.New[string, *exportProcess](),
	}

	// assigning the export engine
	// TODO - need to create multi exporter as per poller count of other point
	if opts.kafka {
		err := exportEngine.setup(string(M.KAFKA), "", opts)
		if err != nil {
			slog.Error("[EXPORT]", "err", err)
			return err
		}
	}
	if opts.icmp {
		err := exportEngine.setup(string(M.DBDUMP), fmt.Sprintf("%s/%d_stat", M.ICMP, M.INTERVAL_60), opts)
		if err != nil {
			slog.Error("[EXPORT] icmp stat", "err", err)
			return err
		}
		err = exportEngine.setup(string(M.DBDUMP), fmt.Sprintf("%s/%d_stat", M.ICMP, M.INTERVAL_300), opts)
		if err != nil {
			slog.Error("[EXPORT] icmp stat", "err", err)
			return err
		}
	}
	// add a webhook exportProcess if needed

	ctx.WithExport(exportEngine)
	return
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

func (e *ExportEngine) setup(name, purpose string, opts Opts) (err error) {
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
		e.ctx, exporterName, db, opts,
	)
	if err != nil {
		return err
	}

	e.process.Set(exporterName, exprt)
	return nil
}
