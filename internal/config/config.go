package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type RemoteTarget struct {
	Host      string   `mapstructure:"host"`
	User      string   `mapstructure:"user"`
	Port      int      `mapstructure:"port"`
	Processes []string `mapstructure:"processes"`
}

type MonitorConfig struct {
	Local  bool           `mapstructure:"local"`
	Remote []RemoteTarget `mapstructure:"remote"`
}

type Config struct {
	Interval string        `mapstructure:"interval"`
	Monitor  MonitorConfig `mapstructure:"monitor"`
}

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

	// Optional: validate interval format
	if _, err := time.ParseDuration(config.Interval); err != nil {
		return nil, fmt.Errorf("invalid interval format: %s", config.Interval)
	}

	return &config, nil
}
