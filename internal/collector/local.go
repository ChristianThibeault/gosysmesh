package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemStats struct {
	Timestamp   time.Time
	CPUPercent  float64
	MemUsedMB   float64
	MemTotalMB  float64
	DiskUsedGB  float64
	DiskTotalGB float64
}

func GetSystemStats() (*SystemStats, error) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("CPU error: %w", err)
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("Mem error: %w", err)
	}

	diskStats, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("Disk error: %w", err)
	}

	return &SystemStats{
		Timestamp:   time.Now(),
		CPUPercent:  cpuPercents[0],
		MemUsedMB:   float64(vm.Used) / 1024 / 1024,
		MemTotalMB:  float64(vm.Total) / 1024 / 1024,
		DiskUsedGB:  float64(diskStats.Used) / 1024 / 1024 / 1024,
		DiskTotalGB: float64(diskStats.Total) / 1024 / 1024 / 1024,
	}, nil
}

