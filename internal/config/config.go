package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
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
	SSHKey         string              `mapstructure:"ssh_key"` 
	ProxyJump      string         	   `mapstructure:"proxy_jump,omitempty"`
    ProcessFilters ProcessFilterConfig `mapstructure:"process_filters"`
}


// JumpConfig defines the configuration for SSH jump hosts.
type JumpConfig struct {
    Host       string `mapstructure:"host"`
    User       string `mapstructure:"user"`
    Port       int    `mapstructure:"port"`
    SSHKeyPath string `mapstructure:"ssh_key"`
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

	// Validate configuration for security
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// validateConfig validates configuration parameters for security
func validateConfig(config *Config) error {
	// Validate interval
	duration, err := time.ParseDuration(config.Interval)
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}
	if duration < time.Second || duration > 24*time.Hour {
		return fmt.Errorf("interval must be between 1 second and 24 hours")
	}

	// Validate remote targets
	for i, target := range config.Monitor.Remote {
		if err := validateRemoteTarget(&target, i); err != nil {
			return fmt.Errorf("remote target %d validation failed: %w", i, err)
		}
	}

	// Validate local process filters
	if err := validateProcessFilters(&config.Monitor.Local.ProcessFilters); err != nil {
		return fmt.Errorf("local process filters validation failed: %w", err)
	}

	return nil
}

// validateRemoteTarget validates a remote target configuration
func validateRemoteTarget(target *RemoteTarget, index int) error {
	// Validate hostname
	if target.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if err := validateHostname(target.Host); err != nil {
		return fmt.Errorf("invalid host: %w", err)
	}

	// Validate username
	if target.User == "" {
		return fmt.Errorf("user cannot be empty")
	}
	if err := validateUsername(target.User); err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}

	// Validate port
	if target.Port < 1 || target.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	// Validate SSH key path
	if target.SSHKey == "" {
		return fmt.Errorf("SSH key path cannot be empty")
	}
	if err := validateFilePath(target.SSHKey); err != nil {
		return fmt.Errorf("invalid SSH key path: %w", err)
	}

	// Validate proxy jump if provided
	if target.ProxyJump != "" {
		if err := validateHostname(target.ProxyJump); err != nil {
			return fmt.Errorf("invalid proxy jump host: %w", err)
		}
	}

	// Validate process filters
	if err := validateProcessFilters(&target.ProcessFilters); err != nil {
		return fmt.Errorf("process filters validation failed: %w", err)
	}

	return nil
}

// validateProcessFilters validates process filter configuration
func validateProcessFilters(filters *ProcessFilterConfig) error {
	// Validate keywords
	for i, keyword := range filters.Keywords {
		if keyword == "" {
			return fmt.Errorf("keyword %d cannot be empty", i)
		}
		if len(keyword) > 100 {
			return fmt.Errorf("keyword %d too long (max 100 characters)", i)
		}
		// Prevent dangerous patterns in keywords
		if strings.ContainsAny(keyword, ";&|$`\n\r") {
			return fmt.Errorf("keyword %d contains dangerous characters", i)
		}
	}

	// Validate users
	for i, user := range filters.Users {
		if err := validateUsername(user); err != nil {
			return fmt.Errorf("user %d validation failed: %w", i, err)
		}
	}

	return nil
}

// validateHostname validates hostname format
func validateHostname(host string) error {
	if host == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	if len(host) > 253 {
		return fmt.Errorf("hostname too long")
	}

	// Allow hostname format (RFC 1123) or IP address
	hostnameRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	ipRegex := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)

	if !hostnameRegex.MatchString(host) && !ipRegex.MatchString(host) {
		return fmt.Errorf("invalid hostname or IP address format")
	}

	return nil
}

// validateUsername validates username format
func validateUsername(user string) error {
	if user == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(user) > 32 {
		return fmt.Errorf("username too long (max 32 characters)")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(user) {
		return fmt.Errorf("invalid username format (only alphanumeric, underscore, dash allowed)")
	}
	return nil
}

// validateFilePath validates file path format
func validateFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	if len(path) > 4096 {
		return fmt.Errorf("file path too long")
	}

	// Clean the path to prevent traversal attacks
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Prevent null bytes and other dangerous characters
	if strings.ContainsAny(path, "\x00\n\r") {
		return fmt.Errorf("file path contains dangerous characters")
	}

	return nil
}
