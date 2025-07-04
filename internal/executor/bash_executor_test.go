package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBashExecutor_RunCommand_SimpleEcho(t *testing.T) {
	executor := NewBashExecutor()

	output, err := executor.RunCommand("echo hello")

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "hello\n", output, "Expected 'hello\\n'")
}

func TestBashExecutor_RunCommand_NonExistentCommand(t *testing.T) {
	executor := NewBashExecutor()

	output, err := executor.RunCommand("nonexistentcommand12345")

	assert.Error(t, err, "Expected error for non-existent command")
	assert.Contains(t, err.Error(), "exit status", "Expected error message to contain 'exit status'")
	assert.True(t, strings.Contains(output, "not found") || strings.Contains(output, "command not found"), "Expected output to contain 'not found' or 'command not found'")
}

func TestBashExecutor_RunCommand_MultiLineOutput(t *testing.T) {
	executor := NewBashExecutor()

	// Create a temporary directory with known structure for testing
	tempDir := t.TempDir()

	// Create some test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tempDir, file))
		require.NoError(t, err, "Failed to create test file")
		f.Close()
	}

	// List files in the temp directory
	output, err := executor.RunCommand("ls " + tempDir)

	require.NoError(t, err, "Expected no error")
	// Check that all test files are in the output
	for _, file := range testFiles {
		assert.Contains(t, output, file, "Expected output to contain '%s'", file)
	}

	// Check that output contains multiple lines (newlines)
	lines := strings.Split(output, "\n")
	assert.GreaterOrEqual(t, len(lines), 3, "Expected at least 3 lines of output")
}

func TestBashExecutor_RunCommand_WithArguments(t *testing.T) {
	executor := NewBashExecutor()

	output, err := executor.RunCommand("echo 'hello world' | wc -w")

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "2", strings.TrimSpace(output), "Expected '2'")
}

func TestBashExecutor_RunCommand_EmptyCommand(t *testing.T) {
	executor := NewBashExecutor()

	output, err := executor.RunCommand("")

	require.NoError(t, err, "Expected no error for empty command")
	assert.Equal(t, "", output, "Expected empty output")
}

func TestBashExecutor_RunCommand_ShellFeatures(t *testing.T) {
	executor := NewBashExecutor()

	// Test that shell features like wildcards work
	// Use a command that relies on shell expansion
	output, err := executor.RunCommand("echo $HOME")

	require.NoError(t, err, "Expected no error")
	assert.NotEqual(t, "$HOME", output, "Expected shell variable expansion, but got literal '$HOME'")
	assert.NotEqual(t, "", output, "Expected non-empty output for $HOME expansion")
}
