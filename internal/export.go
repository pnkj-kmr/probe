package probe

import (
	"fmt"
	"probe/exporter"
	M "probe/model"
	"strconv"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type ExportEngine struct {
	// workers  int
	engine   *actor.Engine
	exporter *safemap.SafeMap[int, *export]
}

func NewExportEngine(workers int, opts ...OptFunc) (*ExportEngine, error) {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return nil, err
	}
	engine := &ExportEngine{
		engine:   e,
		exporter: safemap.New[int, *export](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}
	if options.icmp || options.snmp {
		exprt, err := newExport(KAFKA, "kafka", e, workers)
		if err != nil {
			return nil, err
		}
		engine.exporter.Set(KAFKA, exprt)
	}
	// add a webhook export if needed

	return engine, nil
}

func (e *ExportEngine) Get(id int) (*export, bool) {
	return e.exporter.Get(id)
}

type export struct {
	id      int
	name    string
	workers int
	engine  *actor.Engine
	pid     *actor.PID
	db      *exporter.Exporter
}

func newExport(id int, name string, e *actor.Engine, workers int) (*export, error) {
	if id == KAFKA {
		db, err := exporter.New(id, name)
		if err != nil {
			return nil, err
		}
		return &export{id: id, name: name, engine: e, workers: workers, db: db}, nil
	}
	return &export{id: id, name: name, engine: e, workers: workers}, nil
}

func (s *export) Produce(data []byte) error {
	// on init we need to spin engine spawn with max node
	// here we put mainporc pid to send
	fmt.Println("export produce message received --- ", s.name, data)

	// TODO - need to handle this properly
	// s.pid inbox size to be increated to avoid extra load which found
	//
	s.engine.Send(s.pid, M.ExportMsg{Id: s.id, Name: s.name, Data: data})

	return nil
}

func (s *export) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		fmt.Println("Initialized export --- ", s.name)
	case actor.Started:
		fmt.Println("export started.........", s.name)
	case actor.Stopped:
		fmt.Println("export stopped!!!!!!!!", s.name)
	case M.ExportMsg:
		fmt.Println("========> export message received", msg)
		switch d := msg.Data.(type) {
		case []byte:
			s.db.Export() <- d
		}
	default:
		// need to resend the process
		fmt.Println("message getting exported....", s.name, msg)
		time.Sleep(1 * time.Second)
	}
}

func (s *export) Start() {
	fmt.Println("export---starting------", s.name)
	s.pid = s.engine.SpawnFunc(s.Receive, s.name, actor.WithID(strconv.Itoa(s.id)))

}

func (s *export) Stop() error {
	fmt.Println("export---stopping------", s.name)
	ctx := s.engine.Poison(s.pid)
	<-ctx.Done()
	fmt.Println("export---stopped------", s.name)
	return nil
}
