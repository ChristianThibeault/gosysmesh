package config

import (
    "path/filepath"
    "testing"
)

func TestEdgeCaseConfigs(t *testing.T) {
	basePath := "../../test_configs" 

    tests := []struct {
        name     string
        filename string
        wantErr  bool
    }{
        {"MissingInterval", "missing_interval.yaml", true},
        {"MissingMonitorLocal", "missing_monitor_local.yaml", true},
        {"UnquotedInterval", "unquoted_interval.yaml", true},
        {"InvalidYAML", "invalid_yaml.yaml", true},
        {"NonExistentFile", "does_not_exist.yaml", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            path := filepath.Join(basePath, tt.filename)
            config, err := LoadConfig(path)

            if tt.wantErr && err == nil {
                t.Errorf("Expected error for %s, but got none (config: %+v)", tt.filename, config)
            } else if !tt.wantErr && err != nil {
                t.Errorf("Did not expect error for %s, but got: %v", tt.filename, err)
            }
        })
    }
}
