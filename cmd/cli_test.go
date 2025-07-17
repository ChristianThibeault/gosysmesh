package cmd

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	cmd := exec.Command("../gosysmesh", "--help")
	out, err := cmd.CombinedOutput()

	assert.NoError(t, err)
	assert.Contains(t, string(out), "gosysmesh")
}

