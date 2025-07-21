package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	// Test the help command using the cobra command directly
	var buf bytes.Buffer
	rootCmd.SetOutput(&buf)
	rootCmd.SetArgs([]string{"--help"})
	
	err := rootCmd.Execute()
	assert.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "gosysmesh")
}

