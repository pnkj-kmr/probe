package probe

import (
	"fmt"
	"probe/exporter"
	M "probe/model"
	"strconv"
	"time"

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
	// on init we need to spin engine spawn with max node
	// here we put mainporc pid to send
	fmt.Println("exportProcess produce message received --- ", p.name, data)

	// TODO - need to handle this properly
	// p.pid inbox size to be increated to avoid extra load which found
	//
	p.engine.Send(p.pid, M.ExportMsg{Id: p.id, Name: p.name, Data: data})

	return nil
}

func (p *exportProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		fmt.Println("Initialized exportProcess --- ", p.name)
	case actor.Started:
		fmt.Println("exportProcess started.........", p.name)
	case actor.Stopped:
		fmt.Println("exportProcess stopped!!!!!!!!", p.name)
	case M.ExportMsg:
		fmt.Println("========> exportProcess message received", msg)
		switch d := msg.Data.(type) {
		case []byte:
			p.db.Send() <- d
		}
	default:
		// need to resend the process
		fmt.Println("message getting exported....", p.name, msg)
		time.Sleep(1 * time.Second)
	}
}

func (p *exportProcess) Start() {
	fmt.Println("exportProcess---starting------", p.name)
	p.pid = p.engine.SpawnFunc(p.Receive, p.name, actor.WithID(strconv.Itoa(p.id)))

}

func (p *exportProcess) Stop() error {
	fmt.Println("exportProcess---stopping------", p.name)
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	fmt.Println("exportProcess---stopped------", p.name)
	return nil
}
