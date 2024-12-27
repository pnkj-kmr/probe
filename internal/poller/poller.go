package poller

type Poller interface {
	Poll() error
}

type Creator interface {
	Create(string, []byte) error
}

type Deleter interface {
	Delete(string) error
}

func PINGPoller() {
	// need to sping the ping poller here
	// TODO we should be able to restart poller any time as need
	// with the help to context package

}
