package probe

import (
	"log"
	"log/slog"
	api "probe/apiserver"
	M "probe/model"
	"probe/repository"
	"probe/safemap"

	"github.com/anthdm/hollywood/actor"
)

type apiProcess struct {
	name   M.ProbeType
	engine *actor.Engine
	pid    *actor.PID
	server *api.Server
	db     *DBEngine
}

func newApiProcess(name M.ProbeType, e *actor.Engine, db *DBEngine) (*apiProcess, error) {
	proc := &apiProcess{name: name, engine: e, db: db}
	return proc, nil
}

func (p *apiProcess) Send(data any) error {
	slog.Info("[API] produce a message")
	p.engine.Send(p.pid, data)
	return nil
}

func (p *apiProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		if p.name == M.API {
			var db = safemap.New[string, M.DB]()
			p.db.db.ForEach(func(i string, r *repository.Repository) {
				db.Set(i, r)
			})
			p.server = api.New(p, db)
		}
		slog.Info("[API] process initialized...", "name", p.name)
	case actor.Started:
		// running API chi server here
		go p.server.Run()
		slog.Info("[API] process started", "name", p.name)
	case actor.Stopped:
		slog.Info("[API] process stopped", "name", p.name)

	// need to implement the api heartbeat to make thread up
	// and running for api server
	default:
		_ = msg
		log.Println("[API] default message")
	}
}

func (p *apiProcess) start() {
	slog.Info("[API] process starting...", "name", p.name)
	p.pid = p.engine.SpawnFunc(p.Receive, string(p.name), actor.WithID(""))
}

func (p *apiProcess) stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	slog.Info("[API] process stopped", "name", p.name)
	return nil
}
