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
	Host        string
	Timestamp   time.Time
	Processes   []collector.MonitoredProcess
	SystemStats *collector.SystemStats
}

// CollectRemoteStats collects process info and system stats from a remote server via OpenSSH
func CollectRemoteStats(target config.RemoteTarget) (*RemoteMetrics, error) {
	// Collect process info
	cmd := fmt.Sprintf(`ps -u %s -o pid,user,%%cpu,%%mem,stat,lstart,args --no-headers`, target.User)
	output, err := RunSSHCommandOpenSSH(target.User, target.Host, target.Port, target.SSHKey, target.ProxyJump, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run remote ps on %s: %w", target.Host, err)
	}

	procs := parseProcessOutput(output, target.ProcessFilters)

	// Collect system stats
	systemStats, err := collectRemoteSystemStats(target)
	if err != nil {
		return nil, fmt.Errorf("failed to collect system stats from %s: %w", target.Host, err)
	}

	return &RemoteMetrics{
		Host:        target.Host,
		Timestamp:   time.Now(),
		Processes:   procs,
		SystemStats: systemStats,
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

// collectRemoteSystemStats collects system stats from a remote server via SSH
func collectRemoteSystemStats(target config.RemoteTarget) (*collector.SystemStats, error) {
	// Run multiple commands to get system stats
	cmd := `top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/%us,//'; free -m | awk 'NR==2{printf "%.0f %.0f", $3,$2}'; df -h / | awk 'NR==2{gsub(/[^0-9.]/, "", $3); gsub(/[^0-9.]/, "", $2); printf " %.1f %.1f", $3, $2}'`
	
	output, err := RunSSHCommandOpenSSH(target.User, target.Host, target.Port, target.SSHKey, target.ProxyJump, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run system stats command: %w", err)
	}

	return parseSystemStatsOutput(output)
}

// parseSystemStatsOutput parses system stats from remote command output
func parseSystemStatsOutput(output string) (*collector.SystemStats, error) {
	parts := strings.Fields(strings.TrimSpace(output))
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid system stats output: %s", output)
	}

	cpuPercent, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CPU: %w", err)
	}

	memUsed, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse memory used: %w", err)
	}

	memTotal, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse memory total: %w", err)
	}

	diskUsed, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse disk used: %w", err)
	}

	diskTotal, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse disk total: %w", err)
	}

	return &collector.SystemStats{
		Timestamp:   time.Now(),
		CPUPercent:  cpuPercent,
		MemUsedMB:   memUsed,
		MemTotalMB:  memTotal,
		DiskUsedGB:  diskUsed,
		DiskTotalGB: diskTotal,
	}, nil
}

