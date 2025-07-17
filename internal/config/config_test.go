package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadValidConfig(t *testing.T) {
	yaml := `
interval: "5s"
monitor:
  local:
    enabled: true
    process_filters:
      keywords: ["go"]
      users: ["test"]
      groups: []
  remote: []
`
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yaml))
	assert.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "5s", cfg.Interval)
	assert.True(t, cfg.Monitor.Local.Enabled)
	assert.Contains(t, cfg.Monitor.Local.ProcessFilters.Keywords, "go")
}

func TestLoadInvalidConfig(t *testing.T) {
	yaml := `interval: not-a-duration`
	tmpFile, _ := os.CreateTemp("", "bad-config-*.yaml")
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(yaml))
	tmpFile.Close()

	_, err := LoadConfig(tmpFile.Name())
	assert.Error(t, err)
}

