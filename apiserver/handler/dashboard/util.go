package dashboard

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func getSystemStats() SystemStat {

	// CPU
	cpuPercent, _ := cpu.Percent(0, false) // false = aggregate across all cores
	cpuCores, _ := cpu.Counts(true)        // logical cores
	// fmt.Printf("CPU Usage: %.2f%% (%d cores)\n", cpuPercent[0], cpuCores)

	// Memory
	vmStat, _ := mem.VirtualMemory()
	// fmt.Printf("Memory Total: %.2f GB\n", float64(vmStat.Total)/1e9)
	// fmt.Printf("Memory Used: %.2f GB\n", float64(vmStat.Used)/1e9)
	// fmt.Printf("Memory Usage: %.2f%%\n", vmStat.UsedPercent)

	// Disk
	diskStat, _ := disk.Usage("/")
	// fmt.Printf("Disk Total: %.2f GB\n", float64(diskStat.Total)/1e9)
	// fmt.Printf("Disk Used: %.2f GB\n", float64(diskStat.Used)/1e9)
	// fmt.Printf("Disk Usage: %.2f%%\n", diskStat.UsedPercent)

	return SystemStat{
		CPU:    CPU{Usage: cpuPercent, Cores: cpuCores},
		Memory: *vmStat,
		Disk:   *diskStat,
	}
}
