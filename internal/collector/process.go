package collector

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/ChristianThibeault/gosysmesh/internal/config"
)

// MonitoredProcess represents a filtered process with basic details.
type MonitoredProcess struct {
	PID     int32
	User    string
	Cmdline string
	Name    string
}

// GetFilteredProcesses retrieves processes based on the provided filters.
func GetFilteredProcesses(filters config.ProcessFilterConfig) ([]MonitoredProcess, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	var matches []MonitoredProcess

	for _, p := range procs {
		name, _ := p.Name()
		cmdline, _ := p.Cmdline()
		username, _ := p.Username()

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

