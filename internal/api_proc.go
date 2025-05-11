package probe

import (
	"log"
	"log/slog"
	api "probe/apiserver"
	M "probe/model"
	"probe/repository"
	"probe/safemap"
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
	slog.Info("[API] produce a message")
	p.engine.Send(p.pid, data)
	return nil
}

func (p *apiProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		if p.id == M.API {
			var db = safemap.New[int, M.DB]()
			p.db.db.ForEach(func(i int, r *repository.Repository) {
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
	// case *M.ICMPReq:
	// 	dbId := M.ICMP + msg.Params.PollPeriod
	// 	db, ok := p.db.GetProcess(dbId)
	// 	if ok {
	// 		db.Send(msg)
	// 	} else {
	// 		// this condition will not occurs
	// 		// unless poll_period apart from 60 and 300 seconds
	// 		db, ok := p.db.GetProcess(M.INTERVAL_300 + M.ICMP)
	// 		if ok {
	// 			db.Send(msg)
	// 		}
	// 	}
	default:
		_ = msg
		log.Println("[API] default message")
	}
}

func (p *apiProcess) start() {
	slog.Info("[API] process starting...", "name", p.name)
	p.pid = p.engine.SpawnFunc(p.Receive, p.name, actor.WithID(strconv.Itoa(p.id)))
}

func (p *apiProcess) stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	slog.Info("[API] process stopped", "name", p.name)
	return nil
}
