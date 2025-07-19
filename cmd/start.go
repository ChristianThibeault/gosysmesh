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
					fmt.Printf("Matched processes:\n")
					// for _, p := range filtered {
					// 	fmt.Printf("PID %d: %s (%s) — %s\n", p.PID, p.Name, p.User, p.Cmdline)
					// }
					fmt.Printf("%-25s  %-8s  %-8s  %-6s  %6s  %6s  %6s  %-10s\n", "COMMAND", "TIME", "START", "STAT", "%CPU", "MEM(MB)", "PID", "USER")
					for _, p := range filtered {
						fmt.Printf("%-25s  %-8s  %-8s  %-6s  %6.1f  %6.0f  %6d  %-10s\n",
						p.Name,
						time.Now().Format("15:04:05"),
						p.StartTime,
						p.Status,
						p.CPU,
						p.MEM,
						p.PID,
						p.User,
					)
				}

				}

				for _, target := range conf.Monitor.Remote {
					metrics, err := remote.CollectRemoteStats(target)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Remote %s error: %v\n", target.Host, err)
						continue
					}

					fmt.Printf("[%s][%s] %d processes matched\n",
						metrics.Timestamp.Format("15:04:05"), metrics.Host, len(metrics.Processes),
					)
					fmt.Printf("COMMAND                   TIME      START               STAT   %%CPU  MEM(MB)    PID  USER\n")
					for _, p := range metrics.Processes {
						fmt.Printf("%-25s  %-8s  %-18s  %-6s  %6.1f  %6.0f  %6d  %-10s\n",
						p.Name,
						metrics.Timestamp.Format("15:04:05"),
						p.StartTime,
						p.Status,
						p.CPU,
						p.MEM,
						p.PID,
						p.User,
					)
				}

				}


			case <-quit:
				fmt.Println("Exiting system monitor.")
				break LOOP
			}
		}
	},
}

