package shell

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintHelp_Unit(t *testing.T) {
	var sb strings.Builder
	printHelp(&sb)
	output := sb.String()
	assert.Contains(t, output, "cd <dir>")
	assert.Contains(t, output, "exit")
	assert.Contains(t, output, "help")
	assert.Contains(t, output, "All other input is executed")
}

func TestREPL_HelpIntegration(t *testing.T) {
	binPath := "../binks"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	_ = buildCmd.Run() // Ignore error if already built

	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader("help\nexit\n")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Expected no error from REPL")
	outputStr := string(output)
	assert.Contains(t, outputStr, "Built-in commands:")
	assert.Contains(t, outputStr, "cd <dir>")
	assert.Contains(t, outputStr, "exit")
	assert.Contains(t, outputStr, "help")
	assert.Contains(t, outputStr, "All other input is executed")

	// Test '?' alias
	cmd = exec.Command(binPath)
	cmd.Stdin = strings.NewReader("?\nexit\n")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Expected no error from REPL")
	outputStr = string(output)
	assert.Contains(t, outputStr, "Built-in commands:")
}
