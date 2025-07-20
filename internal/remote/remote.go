package remote

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ChristianThibeault/gosysmesh/internal/collector"
	"github.com/ChristianThibeault/gosysmesh/internal/config"
)

// RemoteMetrics holds data collected from a remote server
type RemoteMetrics struct {
	Host      string
	Timestamp time.Time
	Processes []collector.MonitoredProcess
}

// CollectRemoteStats collects process info from a remote server via OpenSSH
func CollectRemoteStats(target config.RemoteTarget) (*RemoteMetrics, error) {
	cmd := fmt.Sprintf(`ps -u %s -o pid,user,%%cpu,%%mem,stat,lstart,args --no-headers`, target.User)
	output, err := RunSSHCommandOpenSSH(target.User, target.Host, target.Port, target.SSHKey, target.ProxyJump, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run remote ps on %s: %w", target.Host, err)
	}

	procs := parseProcessOutput(output, target.ProcessFilters)

	return &RemoteMetrics{
		Host:      target.Host,
		Timestamp: time.Now(),
		Processes: procs,
	}, nil
}

// parseProcessOutput parses `ps` command output and filters it
func parseProcessOutput(output string, filters config.ProcessFilterConfig) []collector.MonitoredProcess {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var result []collector.MonitoredProcess

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		pid := fields[0]
		user := fields[1]
		cpu := fields[2]
		mem := fields[3]
		stat := fields[4]
		start := strings.Join(fields[5:10], " ") // lstart is 5 fields
		cmdline := strings.Join(fields[10:], " ") // everything after is the command



		proc := collector.MonitoredProcess{
			Name:  cmdline,
			User:  user,
			Group: "",
		}
		if !collector.MatchProcessFilters(proc, filters) {
			continue
		}

		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			continue
		}
		cpuFloat, err := strconv.ParseFloat(cpu, 64)
		if err != nil {
			continue
		}
		memFloat, err := strconv.ParseFloat(mem, 64)
		if err != nil {
			continue
		}

		result = append(result, collector.MonitoredProcess{
			PID:       int32(pidInt),
			User:      user,
			Name:      cmdline,
			Cmdline:   cmdline,
			CPU:       cpuFloat,
			MEM:       memFloat,
			Status:    stat,
			StartTime: start,
		})
	}

	return result
}

