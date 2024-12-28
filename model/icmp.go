package M

import "time"

type ICMPConf struct {
	IP               string        `json:"ip"`
	Cid              string        `json:"cid,omitempty"`
	Metadata         interface{}   `json:"metadata,omitempty"`
	Timeout          time.Duration `json:"timeout,omitempty"`
	Count            int           `json:"count,omitempty"`
	Retries          int           `json:"retries,omitempty"`
	RetryExponential bool          `json:"retry_exponential,omitempty"`
}

type ICMP struct {
	Config     ICMPConf      `json:"config"`
	Err        string        `json:"error,omitempty"`
	Ok         bool          `json:"ok"`
	PacketLoss float64       `json:"packet_loss,omitempty"`
	AvgRtt     time.Duration `json:"avg_rtt,omitempty"`
	StdDevRtt  time.Duration `json:"std_dev_rtt,omitempty"`
}
