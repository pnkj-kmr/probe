package M

import (
	"net/http"
	"time"
)

type PINGParams struct {
	PollPeriod int           `json:"poll_period"`
	PollType   string        `json:"poll_type"`
	PCid       string        `json:"parent_ci_id"`
	MibProfile string        `json:"mib_profile"`
	IP         string        `json:"poll_addr"`
	IPs        []string      `json:"poll_addr_list,omitempty"`
	PktCount   int           `json:"packsiz,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty"`
}

type ICMPReq struct {
	IP         string      `json:"poll_addr"`
	Cid        string      `json:"ci_id"`
	Params     PINGParams  `json:"params,omitempty"`
	InputStats []InputStat `json:"input_stats,omitempty"`
}

type ICMPRes struct {
	Cid    string                `json:"ci_id"`
	Params PINGParams            `json:"params,omitempty"`
	T      string                `json:"t,omitempty"`
	Stats  map[string][4]float64 `json:"stats"`
	Err    string                `json:"error"`
}

func (a *ICMPReq) Bind(r *http.Request) error {
	return nil
}
