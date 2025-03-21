package probe

import (
	"fmt"
	M "probe/model"
	"probe/poller"
	"probe/poller/icmp"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type PollEngine struct {
	poller *safemap.SafeMap[int, *process]
}

func NewPollEngine(db *DBEngine, export *ExportEngine, opts ...OptFunc) (*PollEngine, error) {
	engine := &PollEngine{
		poller: safemap.New[int, *process](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	if options.icmp {
		c, ok := db.Get(ICMP)
		if !ok {
			return nil, ErrDB
		}
		export, ok := export.Get(KAFKA)
		if !ok {
			return nil, ErrNoKAFKA
		}
		engine.poller.Set(ICMP, newProcess(ICMP, "icmp", c, options.workers, export))
	}
	if options.snmp {
		c, ok := db.Get(SNMP)
		if !ok {
			return nil, ErrDB
		}
		export, ok := export.Get(KAFKA)
		if !ok {
			return nil, ErrNoKAFKA
		}
		engine.poller.Set(SNMP, newProcess(SNMP, "snmp", c, options.workers, export))
	}
	return engine, nil
}

func (e *PollEngine) Get(id int) (*process, bool) {
	return e.poller.Get(id)
}

type process struct {
	id       int
	name     string
	db       M.Consumer[<-chan []byte]
	exporter M.Producer[[]byte]
	poller   M.Poller
	workers  int
}

func newProcess(id int, name string, db M.Consumer[<-chan []byte], workers int, exporter M.Producer[[]byte]) *process {
	return &process{id: id, name: name, db: db, workers: workers, exporter: exporter}
}

func (i *process) setPoller() {
	switch i.id {
	case ICMP:
		i.poller = icmp.NewPoller(i.db, i.exporter, poller.WithWorkers(i.workers))
		// case SNMP:
		// i.poller = snmp.NewPoller
	}
}

// Receive - here it's invoke the polling batch on given interval
func (i *process) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		i.setPoller()
		fmt.Println("Initialized poller --- ", i.name)
	case actor.Started:
		fmt.Println("poller started.........", i.name)
	case actor.Stopped:
		fmt.Println("poller stopped!!!!!!!!", i.name)
	case M.PollingBeat:
		fmt.Println("invoking polling again ....", i.id, i.name, msg.Name, ctx.PID().ID)
		i.poller.Poll()
		// for i := 0; i < 10; i++ {
		// 	ctx.Send(ctx.PID(), "hello")
		// 	fmt.Println("msg sent----")
		// }
		// ctx.Send(ctx.PID(), "hello")
	default:
		fmt.Println("message getting computted....", i.name, msg)
		time.Sleep(1 * time.Second)
	}

}

func (i *process) Consume() <-chan []byte {
	return i.db.Consume()
}

// func (i *process) Produce(data []byte) error {
// 	return i.exporter.Produce(data)
// }
