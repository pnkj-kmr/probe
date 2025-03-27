package probe

import (
	"log"
	M "probe/model"
	"probe/poller"
	"probe/poller/icmp"

	"github.com/anthdm/hollywood/actor"
)

type pollProcess struct {
	id       int
	name     string
	db       M.Receiver[<-chan []byte]
	exporter M.Sender[any]
	poller   M.Poller
	workers  int
}

func newPollProcess(id int, name string, db M.Receiver[<-chan []byte], workers int, exporter M.Sender[any]) *pollProcess {
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
		log.Println("[POLL] process initialized...", p.name)
		p.setPoller()
	case actor.Started:
		log.Println("[POLL] process started", p.name)
	case actor.Stopped:
		log.Println("[POLL] process stopped", p.name)
	case M.PollingBeat:
		log.Println("invoking polling again ....", p.id, p.name, msg.Name, ctx.PID().ID)
		p.exporter.Send(M.Do{})
		p.poller.Poll()
		p.exporter.Send(M.Done{})
	default:
		log.Println("[POLL] default poll process message")
		_ = msg
	}

}
