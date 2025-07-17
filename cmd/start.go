package cmd

import (
	"fmt"
	"os"

	"github.com/ChristianThibeault/gosysmesh/internal/config"
	"github.com/spf13/cobra"
)

// var cfgFile string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start system monitoring based on config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Loaded config:")
		fmt.Printf("Interval: %s\n", conf.Interval)

		if conf.Monitor.Local {
			fmt.Println("✅ Local monitoring enabled.")
			// Start local collector here...
		}

		for _, remote := range conf.Monitor.Remote {
			fmt.Printf("🌐 Remote: %s@%s (port %d), tracking: %v\n",
				remote.User, remote.Host, remote.Port, remote.Processes)
			// Connect via SSH and collect...
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&cfgFile, "config", "config.yaml", "Path to config file")
}

