package probe

import (
	"fmt"
	M "probe/model"
	"probe/poller"
	"probe/poller/icmp"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type pollProcess struct {
	id       int
	name     string
	db       M.Receiver[<-chan []byte]
	exporter M.Sender[[]byte]
	poller   M.Poller
	workers  int
}

func newPollProcess(id int, name string, db M.Receiver[<-chan []byte], workers int, exporter M.Sender[[]byte]) *pollProcess {
	return &pollProcess{id: id, name: name, db: db, workers: workers, exporter: exporter}
}

func (p *pollProcess) setPoller() {
	switch p.id {
	case ICMP:
		p.poller = icmp.NewPoller(p.db, p.exporter, poller.WithWorkers(p.workers))
		// case SNMP:
		// p.poller = snmp.NewPoller
	}
}

// Receive - here it's invoke the polling batch on given interval
func (p *pollProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		p.setPoller()
		fmt.Println("Initialized poller --- ", p.name)
	case actor.Started:
		fmt.Println("poller started.........", p.name)
	case actor.Stopped:
		fmt.Println("poller stopped!!!!!!!!", p.name)
	case M.PollingBeat:
		fmt.Println("invoking polling again ....", p.id, p.name, msg.Name, ctx.PID().ID)
		p.poller.Poll()
		// for p := 0; p < 10; p++ {
		// 	ctx.Send(ctx.PID(), "hello")
		// 	fmt.Println("msg sent----")
		// }
		// ctx.Send(ctx.PID(), "hello")
	default:
		fmt.Println("message getting computted....", p.name, msg)
		time.Sleep(1 * time.Second)
	}

}

// func (p *pollProcess) Consume() <-chan []byte {
// 	return p.db.Receive()
// }

// func (p *pollProcess) Produce(data []byte) error {
// 	return p.exporter.Produce(data)
// }
