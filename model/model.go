package M

import (
	"encoding/json"
)

type None struct{}

// process start/end aknowledge types
type Do struct{}
type Done struct{}

type Record struct {
	Key   []byte
	Value []byte
}

type PollingBeat struct {
	Id   int
	Name string
}

type ExportMsg struct {
	Data        any
	Agent       string
	Timezone    int32
	PartitionId int32
}

func (x ExportMsg) Encode() ([]byte, error) {
	// var b bytes.Buffer
	// enc := gob.NewEncoder(&b)
	// err := enc.Encode(x)
	// if err != nil {
	// 	log.Println("ENCODE ERROR", err)
	// }
	// log.Println("message ----", b.Bytes())
	// return b.Bytes(), err
	return json.Marshal(x)
}

func (x ExportMsg) Length() int {
	// var b bytes.Buffer
	// enc := gob.NewEncoder(&b)
	// enc.Encode(x)
	// return len(b.Bytes())
	b, _ := json.Marshal(x)
	return len(b)
}
