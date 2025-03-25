package probe

import (
	"fmt"
	api "probe/apiserver"
	M "probe/model"
	"strconv"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type apiProcess struct {
	id     int
	name   string
	engine *actor.Engine
	pid    *actor.PID
	server *api.Server
}

func newApiProcess(id int, name string, e *actor.Engine) (*apiProcess, error) {
	proc := &apiProcess{id: id, name: name, engine: e}
	if id == API {
		proc.server = api.New(proc)
	}
	return proc, nil
}

func (p *apiProcess) Send(data []byte) error {
	// on init we need to spin engine spawn with max node
	// here we put mainporc pid to send
	fmt.Println("apiProcess produce message received --- ", p.name, data)

	// TODO - need to handle this properly
	// p.pid inbox size to be increated to avoid extra load which found
	//
	p.engine.Send(p.pid, M.ExportMsg{Id: p.id, Name: p.name, Data: data})

	return nil
}

func (p *apiProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		fmt.Println("Initialized apiProcess --- ", p.name)
	case actor.Started:
		fmt.Println("apiProcess started.........", p.name)
	case actor.Stopped:
		fmt.Println("apiProcess stopped!!!!!!!!", p.name)

	default:
		// need to resend the process
		fmt.Println("message getting... ,,,, from api to save into system...", p.name, msg)
		time.Sleep(1 * time.Second)
	}
}

func (p *apiProcess) start() {
	fmt.Println("apiProcess---starting------", p.name)
	p.pid = p.engine.SpawnFunc(p.Receive, p.name, actor.WithID(strconv.Itoa(p.id)))
	// running API chi server here
	go p.server.Run()
}

func (p *apiProcess) stop() error {
	fmt.Println("apiProcess---stopping------", p.name)
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	fmt.Println("apiProcess---stopped------", p.name)
	return nil
}
