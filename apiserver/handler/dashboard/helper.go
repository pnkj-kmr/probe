package dashboard

import (
	"fmt"
	M "probe/model"
)

func (api *R) resourceCount(protocol M.ProbeType) (t uint64) {
	p60 := fmt.Sprintf("%s/%d", protocol, M.INTERVAL_60)
	if db, ok := api.DB.Get(p60); ok {
		total := db.Count()
		t = t + total
	}

	p300 := fmt.Sprintf("%s/%d", protocol, M.INTERVAL_300)
	if db, ok := api.DB.Get(p300); ok {
		total := db.Count()
		t = t + total
	}

	return
}
