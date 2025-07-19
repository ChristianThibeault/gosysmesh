package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/ChristianThibeault/gosysmesh/internal/config"
)

// MonitoredProcess represents a filtered process with basic details.
type MonitoredProcess struct {
	PID       int32
	User      string
	Group     string
	Name      string
	Cmdline   string
	CPU       float64
	MEM       float64
	StartTime string 
	Status    string 
}

// GetFilteredProcesses retrieves processes based on the provided filters.
func GetFilteredProcesses(filters config.ProcessFilterConfig) ([]MonitoredProcess, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	var matches []MonitoredProcess

	for _, p := range procs {
		// name, _ := p.Name()
		// cmdline, _ := p.Cmdline()
		// username, _ := p.Username()

		name, _ := p.Name()
		cmdline, _ := p.Cmdline()
		username, _ := p.Username()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		startTime, _ := p.CreateTime() // returns Unix milliseconds
		statusList, _ := p.Status()
		status := ""
		if len(statusList) > 0 {
			status = statusList[0]
		}


		start := time.UnixMilli(startTime).Format("15:04:05")


		if !matchesKeyword(name, cmdline, filters.Keywords) {
			continue
		}
		if len(filters.Users) > 0 && !stringInSlice(username, filters.Users) {
			continue
		}

		matches = append(matches, MonitoredProcess{
			PID:     p.Pid,
			User:    username,
			Name:    name,
			Cmdline: cmdline,
			CPU:       cpuPercent,
			MEM:       float64(memPercent),
			Status:    status,
			StartTime: start,
		})



	}

	return matches, nil
}

// matchesKeyword checks if the process name or command line contains any of the specified keywords.
func matchesKeyword(name, cmdline string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(name, kw) || strings.Contains(cmdline, kw) {
			return true
		}
	}
	return false
}

// stringInSlice checks if a string is present in a slice of strings.
func stringInSlice(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

// MatchProcessFilters returns true if the given MonitoredProcess matches the filter
func MatchProcessFilters(proc MonitoredProcess, filters config.ProcessFilterConfig) bool {
    for _, keyword := range filters.Keywords {
        if strings.Contains(proc.Name, keyword) {
            return true
        }
    }
    for _, user := range filters.Users {
        if proc.User == user {
            return true
        }
    }
    for _, group := range filters.Groups {
        if proc.Group == group {
            return true
        }
    }
    return false
}

