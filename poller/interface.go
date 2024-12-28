package poller

type Poller interface {
	Poll() error
}

type Scanner interface {
	Scan(interface{}) (interface{}, error)
}
