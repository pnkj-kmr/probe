package profile

import (
	"encoding/json"
	"fmt"
	M "probe/model"
)

func (api *R) getProcessStatus() (data map[string]any, err error) {
	db, ok := api.DB.Get(string(M.PROCESS))
	if !ok {
		return data, fmt.Errorf("NO_PROCESS_DETAIL")
	}

	var p M.ProcessBeat
	data = make(map[string]any)
	for v := range db.Receive() {
		err := json.Unmarshal(v, &p)
		if err != nil {
			return data, err
		}
		data[p.Name] = p
	}
	return
}
