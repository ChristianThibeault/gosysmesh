package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ChristianThibeault/gosysmesh/internal/config"
)

func TestGetFilteredProcesses(t *testing.T) {
	filters := config.ProcessFilterConfig{
		Keywords: []string{"go"},
		Users:    []string{},
		Groups:   []string{},
	}

	procs, err := GetFilteredProcesses(filters)
	assert.NoError(t, err)
	assert.NotNil(t, procs) // Might be empty, but shouldn’t crash
}

