package M

import (
	"encoding/json"
)

type (
	None struct{}

	// process start/end aknowledge types
	Do   struct{}
	Done struct{}

	Record struct {
		K string
		V any
	}

	PollingBeat struct {
		Name string
	}

	ProcessBeat struct {
		Counter uint64 `json:"counter"`
		Name    string `json:"name"`
		St      string `json:"st"`
		Et      string `json:"et"`
		T       string `json:"t"`
		Given   int    `json:"given"`
		Polled  int    `json:"polled"`
		Error   string `json:"error"`
	}
)
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
