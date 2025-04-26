package probe

import (
	"log"
	api "probe/apiserver"
	M "probe/model"
	"strconv"

	"github.com/anthdm/hollywood/actor"
)

type apiProcess struct {
	id     int
	name   string
	engine *actor.Engine
	pid    *actor.PID
	server *api.Server
	db     *DBEngine
}

func newApiProcess(id int, name string, e *actor.Engine, db *DBEngine) (*apiProcess, error) {
	proc := &apiProcess{id: id, name: name, engine: e, db: db}
	return proc, nil
}

func (p *apiProcess) Send(data any) error {
	log.Println("[API] produce a message")
	p.engine.Send(p.pid, data)
	return nil
}

func (p *apiProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		if p.id == API {
			p.server = api.New(p)
		}
		log.Println("[API] process initialized...", p.name)
	case actor.Started:
		// running API chi server here
		go p.server.Run()
		log.Println("[API] process started", p.name)
	case actor.Stopped:
		log.Println("[API] process stopped", p.name)
	case *M.ICMPReq:
		db, ok := p.db.GetProcess(ICMP)
		if ok {
			db.Send(msg)
		}
	default:
		_ = msg
		log.Println("[API] default message")
	}
}

func (p *apiProcess) start() {
	// log.Println("[API] process starting...")
	p.pid = p.engine.SpawnFunc(p.Receive, p.name, actor.WithID(strconv.Itoa(p.id)))
}

func (p *apiProcess) stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	// log.Println("[API] process stopped")
	return nil
}
