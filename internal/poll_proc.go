package probe

import (
	"log/slog"
	M "probe/model"
	"probe/poller"
	"probe/poller/icmp"
	"probe/poller/snmp"
	"strings"

	"github.com/anthdm/hollywood/actor"
)

type pollProcess struct {
	name     string
	db       M.Receiver[<-chan []byte]
	exporter []M.Sender[any]
	finder   M.Finder
	poller   M.Poller
	workers  int
}

func newPollProcess(name string, db M.Receiver[<-chan []byte], workers int, exporter []M.Sender[any], finder M.Finder) *pollProcess {
	return &pollProcess{name: name, db: db, workers: workers, exporter: exporter, finder: finder}
}

func (p *pollProcess) setPoller() {
	var _pollType M.ProbeType
	if strings.Contains(p.name, "icmp") {
		_pollType = M.ICMP
	} else if strings.Contains(p.name, "snmp") {
		_pollType = M.SNMP
	}
	switch _pollType {
	case M.ICMP:
		p.poller = icmp.NewPoller(p.db, p.exporter, poller.WithWorkers(p.workers))
	case M.SNMP:
		p.poller = snmp.NewPoller(p.db, p.exporter, p.finder, poller.WithWorkers(p.workers))
	}
}

// Receive - here it's invoke the polling batch on given interval
func (p *pollProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		p.setPoller()
		slog.Info("[POLL] process initialized...", "name", p.name)
	case actor.Started:
		slog.Info("[POLL] process started", "name", p.name)
	case actor.Stopped:
		slog.Info("[POLL] process stopped", "name", p.name)
	case M.PollingBeat:
		slog.Info("invoking poller....", "id", ctx.PID().ID, "name", p.name, "msg", msg.Name, "exporter", p.exporter)
		// p.exporter.Send(M.Do{})
		for _, sender := range p.exporter {
			sender.Send(M.Do{})
		}
		if p.poller != nil {
			p.poller.Poll()
		} else {
			slog.Warn("No poller found", "poller", p.poller)
		}
		// p.exporter.Send(M.Done{})
		for _, sender := range p.exporter {
			sender.Send(M.Done{})
		}
	default:
		slog.Info("[POLL] default poll process message")
		_ = msg
	}

}
