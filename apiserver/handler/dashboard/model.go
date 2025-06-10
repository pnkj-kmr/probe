package dashboard

type CPU struct {
	Usage []float64 `json:"usage"`
	Cores int       `json:"cores"`
}

type Memory struct {
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type Disk struct {
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type SystemStats struct {
	CPU    CPU    `json:"cpu"`
	Memory Memory `json:"memory"`
	Disk   Disk   `json:"disk"`
}

type PollStat struct {
	Name   string `json:"name"`
	Total  uint64 `json:"total"`
	Polled uint64 `json:"polled"`
}

type ResourceStats struct {
	Stats []PollStat `json:"stats"`
	Total uint64     `json:"total"`
}

type PollingStats []PollStat
