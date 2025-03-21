package M

type None struct{}

type Record struct {
	Key   []byte
	Value []byte
}

type PollingBeat struct {
	Id   int
	Name string
}

type ExportMsg struct {
	Id   int
	Name string
	Data any
}
