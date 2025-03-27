package M

import (
	"net/http"
	"time"
)

type PingParams struct {
	PollPeriod int           `json:"poll_period"` // not in use - default poll interval used
	PollType   string        `json:"poll_type"`
	PCid       string        `json:"parent_ci_id"`
	MibProfile string        `json:"mib_profile"`
	IP         string        `json:"poll_addr"`
	IPs        []string      `json:"ips,omitempty"`
	PktCount   int           `json:"packsiz,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty"`
}

type InICMP struct {
	IP         string      `json:"ip"`
	Cid        string      `json:"ci_id"`
	Params     PingParams  `json:"params,omitempty"`
	InputStats []InputStat `json:"input_stats,omitempty"`
}

type OutICMP struct {
	Cid    string                `json:"ci_id"`
	Params PingParams            `json:"params,omitempty"`
	T      string                `json:"datetime,omitempty"`
	Stats  map[string][4]float64 `json:"stats"`
	Err    string                `json:"error"`
}

////TODO - need to do default packet configuration while initlisating the poller configuration
// input-----
// Metadata         any           `json:"metadata,omitempty"`
// Timeout          time.Duration `json:"timeout,omitempty"`
// Count            int           `json:"count,omitempty"`
// Retries          int           `json:"retries,omitempty"`
// RetryExponential bool          `json:"retry_exponential,omitempty"`
// output-----
// Stats  map[string][4]float64
// Config     InICMP        `json:"config"`
// Err        string        `json:"error,omitempty"`
// Ok         bool          `json:"ok"`
// PacketLoss float64       `json:"packet_loss,omitempty"`
// AvgRtt     time.Duration `json:"avg_rtt,omitempty"`
// StdDevRtt  time.Duration `json:"std_dev_rtt,omitempty"`

func (a *InICMP) Bind(r *http.Request) error {
	return nil
}
