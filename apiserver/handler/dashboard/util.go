package dashboard

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}

func getSystemStats() SystemStats {

	// CPU
	cpuPercent, _ := cpu.Percent(time.Second*2, false) // false = aggregate across all cores
	cpuCores, _ := cpu.Counts(true)                    // logical cores
	// fmt.Printf("CPU Usage: %.2f%% (%d cores)\n", cpuPercent[0], cpuCores)

	// Memory
	vmStat, _ := mem.VirtualMemory()
	// fmt.Printf("Memory Total: %.2f GB\n", float64(vmStat.Total)/1e9)
	// fmt.Printf("Memory Used: %.2f GB\n", float64(vmStat.Used)/1e9)
	// fmt.Printf("Memory Usage: %.2f%%\n", vmStat.UsedPercent)
	memory := Memory{
		Total:       fmt.Sprintf("%.2f GB", float64(vmStat.Total)/1e9),
		Used:        fmt.Sprintf("%.2f GB", float64(vmStat.Used)/1e9),
		UsedPercent: roundToTwoDecimals(vmStat.UsedPercent),
	}

	// Disk
	diskStat, _ := disk.Usage("/")
	// fmt.Printf("Disk Total: %.2f GB\n", float64(diskStat.Total)/1e9)
	// fmt.Printf("Disk Used: %.2f GB\n", float64(diskStat.Used)/1e9)
	// fmt.Printf("Disk Usage: %.2f%%\n", diskStat.UsedPercent)
	disk := Disk{
		Total:       fmt.Sprintf("%.2f GB", float64(diskStat.Total)/1e9),
		Used:        fmt.Sprintf("%.2f GB", float64(diskStat.Used)/1e9),
		UsedPercent: roundToTwoDecimals(diskStat.UsedPercent),
	}

	return SystemStats{
		CPU:    CPU{Usage: roundToTwoDecimals(cpuPercent[0]), Cores: cpuCores},
		Memory: memory,
		Disk:   disk,
	}
}
