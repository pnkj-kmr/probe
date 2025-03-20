package probe

import "github.com/anthdm/hollywood/safemap"

type APIEngine struct {
	db *safemap.SafeMap[int, any]
}

func NewAPIEngine(opts ...OptFunc) (*APIEngine, error) {
	engine := &APIEngine{
		db: safemap.New[int, any](),
	}
	options := DefaultOpts()
	for _, opt := range opts {
		opt(&options)
	}

	return engine, nil
}

func (e *APIEngine) Get(id int) (any, bool) {
	return e.db.Get(id)
}
