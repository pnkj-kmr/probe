package probe

import (
	"log"
	"probe/exporter"
	M "probe/model"
	"strconv"

	"github.com/anthdm/hollywood/actor"
)

type exportProcess struct {
	id      int
	name    string
	workers int
	engine  *actor.Engine
	pid     *actor.PID
	db      *exporter.Exporter // M.Sender - need to make
}

func newExportProcess(id int, name string, e *actor.Engine, workers int) (*exportProcess, error) {
	if id == KAFKA {
		db, err := exporter.New(id, name)
		if err != nil {
			return nil, err
		}
		return &exportProcess{id: id, name: name, engine: e, workers: workers, db: db}, nil
	}
	return &exportProcess{id: id, name: name, engine: e, workers: workers}, nil
}

func (p *exportProcess) Send(data []byte) error {
	p.engine.Send(p.pid, M.ExportMsg{Id: p.id, Name: p.name, Data: data})
	return nil
}

func (p *exportProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		log.Println("[EXPORT] process initialized...", p.name)
	case actor.Started:
		log.Println("[EXPORT] process started", p.name)
	case actor.Stopped:
		log.Println("[EXPORT] process stopped", p.name)
	case M.ExportMsg:
		log.Println("========> exportProcess message received", msg)
		switch d := msg.Data.(type) {
		case []byte:
			p.db.Send() <- d
		}
	default:
		log.Println("[EXPORT] default message")
		_ = msg
	}
}

func (p *exportProcess) Start() {
	// log.Println("[EXPORT] process starting...")
	p.pid = p.engine.SpawnFunc(p.Receive, p.name, actor.WithID(strconv.Itoa(p.id)))

}

func (p *exportProcess) Stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	// log.Println("[EXPORT] process stopped")
	return nil
}
