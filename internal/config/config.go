package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// ProcessFilterConfig defines the filtering criteria for processes.
type ProcessFilterConfig struct {
    Keywords []string `mapstructure:"keywords"`
    Users    []string `mapstructure:"users"`
    Groups   []string `mapstructure:"groups"`
}

// LocalMonitorConfig defines the configuration for local process monitoring.
type LocalMonitorConfig struct {
    Enabled        bool                `mapstructure:"enabled"`
    ProcessFilters ProcessFilterConfig `mapstructure:"process_filters"`
}

// RemoteTarget defines the configuration for remote process monitoring targets.
type RemoteTarget struct {
    Host           string              `mapstructure:"host"`
    User           string              `mapstructure:"user"`
    Port           int                 `mapstructure:"port"`
    ProcessFilters ProcessFilterConfig `mapstructure:"process_filters"`
}

// MonitorConfig aggregates local and remote monitoring configurations.
type MonitorConfig struct {
    Local  LocalMonitorConfig `mapstructure:"local"`
    Remote []RemoteTarget     `mapstructure:"remote"`
}

// Config structure for the gosysmesh application.
type Config struct {
    Interval string        `mapstructure:"interval"`
    Monitor  MonitorConfig `mapstructure:"monitor"`
}

// LoadConfig reads the configuration from a YAML file and unmarshals it into a Config struct.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if !viper.IsSet("monitor.local") {
    return nil, fmt.Errorf("missing required field: monitor.local")
	}



	if _, err := time.ParseDuration(config.Interval); err != nil {
		return nil, fmt.Errorf("invalid interval format: %s", config.Interval)
	}

	return &config, nil
}
