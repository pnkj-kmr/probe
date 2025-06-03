package dashboard

import (
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type CPU struct {
	Usage []float64 `json:"usage"`
	Cores int       `json:"cores"`
}

type SystemStat struct {
	CPU    CPU                   `json:"cpu"`
	Memory mem.VirtualMemoryStat `json:"memory"`
	Disk   disk.UsageStat        `json:"disk"`
}
