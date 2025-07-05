package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func containsPrompt(output string) bool {
	plain := "binks>"
	colored := "binks:"
	return strings.Contains(output, plain) || strings.Contains(output, colored)
}

func TestCLI_NoArgs_StartsREPL(t *testing.T) {
	binPath := "../binks"
	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/binks")
	_ = buildCmd.Run() // Ignore error if already built

	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader("exit\n")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.True(t, containsPrompt(string(output)), "Expected prompt in output")
}

func TestCLI_HelpCommand(t *testing.T) {
	binPath := "../binks"
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader("help\nexit\n")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	outStr := string(output)
	assert.Contains(t, outStr, "Built-in commands:")
	assert.Contains(t, outStr, "cd <dir>")
	assert.Contains(t, outStr, "exit")
	assert.Contains(t, outStr, "help")
}

func TestCLI_CD_Builtin(t *testing.T) {
	binPath := "../binks"
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader("cd /\npwd\nexit\n")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	outStr := string(output)
	assert.Contains(t, outStr, "/\n") // Should print root dir
}

func TestCLI_InvalidCommand_ErrorOutput(t *testing.T) {
	binPath := "../binks"
	cmd := exec.Command(binPath, "invalidcommand12345")
	output, err := cmd.CombinedOutput()
	assert.Error(t, err)
	assert.Contains(t, string(output), "Error:")
}
