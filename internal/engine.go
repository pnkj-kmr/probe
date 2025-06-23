package probe

import (
	"log"

	"github.com/anthdm/hollywood/actor"
)

func newEngine() (*actor.Engine, error) {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Println("[ENGINE] unable to start the engine")
		return nil, err
	}
	return e, err
}
