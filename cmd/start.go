package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChristianThibeault/gosysmesh/internal/collector"
	"github.com/ChristianThibeault/gosysmesh/internal/config"
	"github.com/ChristianThibeault/gosysmesh/internal/remote"
	"github.com/spf13/cobra"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
)


func printHostProcesses(title string, timestamp time.Time, procs []collector.MonitoredProcess) {
	fmt.Printf("%s%s%s%s [%s]%s\n", bold, cyan, title, reset, timestamp.Format("15:04:05"), reset)

	for i, p := range procs {
		conn := "├──"
		if i == len(procs)-1 {
			conn = "└──"
		}

		cpuColor := green
		switch {
		case p.CPU > 70:
			cpuColor = red
		case p.CPU > 30:
			cpuColor = yellow
		}

		memColor := green
		switch {
		case p.MEM > 70: // >70%
			memColor = red
		case p.MEM > 30: // >30%
			memColor = yellow
		}

		fmt.Printf("%s PID %-6d: %s\n", conn, p.PID, p.Cmdline)
		fmt.Printf("│   ├── %sCPU:%s %.1f%%   %sMEM:%s %.1f%%\n", cpuColor, reset, p.CPU, memColor, reset, p.MEM)
		fmt.Printf("│   └── Start: %s   Stat: %s   User: %s%s%s\n",
			p.StartTime, p.Status, blue, p.User, reset)
	}
	fmt.Println()
}


var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start system monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		interval, err := time.ParseDuration(conf.Interval)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid interval: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Starting local system monitor: every %s\n", interval)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		LOOP:
		for {
			select {
			case <-ticker.C:


				stats, err := collector.GetSystemStats()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error collecting stats: %v\n", err)
					continue
				}

				fmt.Printf("[%s] CPU: %.1f%% | MEM: %.0f/%.0f MB | DISK: %.1f/%.1f GB\n",
				stats.Timestamp.Format("15:04:05"),
				stats.CPUPercent,
				stats.MemUsedMB, stats.MemTotalMB,
				stats.DiskUsedGB, stats.DiskTotalGB,
			)

			filtered, err := collector.GetFilteredProcesses(conf.Monitor.Local.ProcessFilters)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error filtering processes: %v\n", err)
			} else {
				printHostProcesses("local", time.Now(), filtered)
			}

			for _, target := range conf.Monitor.Remote {
				metrics, err := remote.CollectRemoteStats(target)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Remote %s error: %v\n", target.Host, err)
					continue
				}

				// Print remote system stats
				if metrics.SystemStats != nil {
					fmt.Printf("[%s][%s] CPU: %.1f%% | MEM: %.0f/%.0f MB | DISK: %.1f/%.1f GB\n",
						metrics.Timestamp.Format("15:04:05"), metrics.Host,
						metrics.SystemStats.CPUPercent,
						metrics.SystemStats.MemUsedMB, metrics.SystemStats.MemTotalMB,
						metrics.SystemStats.DiskUsedGB, metrics.SystemStats.DiskTotalGB,
					)
				}

				if len(metrics.Processes) > 0 {
					printHostProcesses(target.Host, metrics.Timestamp, metrics.Processes)
				}
			}


			case <-quit:
				fmt.Println("Exiting system monitor.")
				break LOOP
			}
		}
	},
}

