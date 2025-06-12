package dashboard

import (
	"fmt"
	M "probe/model"
)

func (api *R) resCount(k string) (t uint64) {
	if db, ok := api.DB.Get(k); ok {
		total := db.Count()
		t = t + total
	}
	return t
}

func (api *R) resourceCount(protocol M.ProbeType) (t uint64) {
	p60 := fmt.Sprintf("%s/%d", protocol, M.INTERVAL_60)
	t = t + api.resCount(p60)

	p300 := fmt.Sprintf("%s/%d", protocol, M.INTERVAL_300)
	t = t + api.resCount(p300)

	return
}
