package probe

import (
	"encoding/json"
	"log"
	M "probe/model"
	"probe/repository"
	"strconv"

	"github.com/anthdm/hollywood/actor"
)

type dbProcess struct {
	id     int
	name   string
	engine *actor.Engine
	pid    *actor.PID
	db     *repository.Repository
}

func newDBProcess(id int, name string, e *actor.Engine, db *repository.Repository) (*dbProcess, error) {
	return &dbProcess{
		id: id, name: name,
		engine: e,
		db:     db,
	}, nil
}

func (p *dbProcess) Send(data any) error {
	p.engine.Send(p.pid, data)
	return nil
}

func (p *dbProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		log.Println("[DB] process initialized...", p.name)
	case actor.Started:
		log.Println("[DB] process started", p.name)
	case actor.Stopped:
		log.Println("[DB] process stopped", p.name)
	case *M.ICMPReq:
		log.Println("==== ping record received ---", msg)
		d, err := json.Marshal(msg)
		if err != nil {
			// todo - need to handle the below case
			log.Println("[ERROR] ============ ", err)
		}
		err = p.db.Create(msg.IP, d)
		log.Println("[DB SAVED] ============ ", d, err)

	default:
		_ = msg
		log.Println("[DB] default message")
	}
}

func (p *dbProcess) start() {
	// log.Println("[DB] process starting...")
	p.pid = p.engine.SpawnFunc(p.Receive, p.name, actor.WithID(strconv.Itoa(p.id)))
}

func (p *dbProcess) stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	// log.Println("[DB] process stopped")
	return nil
}
