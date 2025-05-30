package M

import (
	"time"
)

type PINGParams struct {
	PollPeriod int           `json:"poll_period" validate:"required,oneof=5 60 300"`
	PollType   string        `json:"poll_type" validate:"required,oneof=ping"`
	PCid       string        `json:"parent_ci_id"`
	MibProfile string        `json:"mib_profile"`
	IP         string        `json:"poll_addr"`
	IPs        []string      `json:"poll_addr_list,omitempty"`
	PktCount   int           `json:"packsiz,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty" swaggerignore:"true"`
}

type ICMPReq struct {
	IP         string      `json:"poll_addr" validate:"required"`
	Cid        string      `json:"ci_id" validate:"required"`
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
