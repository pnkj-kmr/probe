package probe

import (
	"encoding/json"
	"fmt"
	"log/slog"
	M "probe/model"

	"github.com/anthdm/hollywood/actor"
)

type eventProcess struct {
	name   string
	db     *DBEngine
	engine *actor.Engine
	pid    *actor.PID
}

func newEventProcess(name string, e *actor.Engine, db *DBEngine) (*eventProcess, error) {
	return &eventProcess{
		name: name, engine: e, db: db,
	}, nil
}

func (p *eventProcess) Start() {
	p.pid = p.engine.SpawnFunc(p.Receive, p.name)
}

func (p *eventProcess) Stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	return nil
}

func (p *eventProcess) Send(data any) error {
	// TODO - any channel need to be handle
	p.engine.Send(p.pid, data)
	return nil
}

func (p *eventProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		slog.Info("[EVENT] process initialized...", "name", p.name)
	case actor.Started:
		slog.Info("[EVENT] process started", "name", p.name)
	case actor.Stopped:
		slog.Info("[EVENT] process stopped", "name", p.name)
	case M.ProcessBeat:
		fmt.Println("need to sve intp db")
		p.saveProcessStatus(msg)
	default:
		slog.Info("[EVENT] default message")
		_ = msg
	}
}

func (p *eventProcess) saveProcessStatus(data M.ProcessBeat) (err error) {
	db, ok := p.db.GetDB(string(M.PROCESS))
	if !ok {
		return fmt.Errorf("NO_PROCESS_DB")
	}

	ex, err := db.Find(data.Name)
	if err == nil {
		var exData M.ProcessBeat
		json.Unmarshal(ex, &exData)
		data.Counter = exData.Counter + 1
	}

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = db.Create(data.Name, d)
	return
}
