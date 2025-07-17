package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChristianThibeault/gosysmesh/internal/collector"
	"github.com/ChristianThibeault/gosysmesh/internal/config"
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
			case <-quit:
				fmt.Println("Exiting system monitor.")
				break LOOP
			}
		}
	},
}

