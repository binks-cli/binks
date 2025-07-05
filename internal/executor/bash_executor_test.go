package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBashExecutor_RunCommand_TableDriven(t *testing.T) {
	executor := NewBashExecutor()
	testCases := []struct {
		name        string
		command     string
		expect      string
		expectError bool
		outputCheck func(_ *testing.T, output string)
		errCheck    func(_ *testing.T, err error)
	}{
		{
			name:    "simple echo",
			command: "echo hello",
			expect:  "hello\n",
		},
		{
			name:        "non-existent command",
			command:     "nonexistentcommand12345",
			expectError: true,
			outputCheck: func(t *testing.T, output string) {
				assert.True(t, strings.Contains(output, "not found") || strings.Contains(output, "command not found"), "Expected output to contain 'not found' or 'command not found'")
			},
			errCheck: func(t *testing.T, err error) {
				assert.Error(t, err, "Expected error for non-existent command")
				assert.Contains(t, err.Error(), "exit status", "Expected error message to contain 'exit status'")
			},
		},
		{
			name:    "with arguments",
			command: "echo 'hello world' | wc -w",
			outputCheck: func(t *testing.T, output string) {
				// Accept both '2' and '       2' (with or without leading spaces)
				trimmed := strings.TrimSpace(output)
				assert.Equal(t, "2", trimmed, "Expected word count to be '2'")
			},
		},
		{
			name:    "empty command",
			command: "",
			expect:  "",
		},
		{
			name:    "shell features (env expansion)",
			command: "echo $HOME",
			outputCheck: func(t *testing.T, output string) {
				assert.NotEqual(t, "$HOME", output, "Expected shell variable expansion, but got literal '$HOME'")
				assert.NotEqual(t, "", output, "Expected non-empty output for $HOME expansion")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := executor.RunCommand(tc.command)
			if tc.expectError {
				if tc.errCheck != nil {
					tc.errCheck(t, err)
				} else {
					assert.Error(t, err)
				}
			} else {
				require.NoError(t, err, "Expected no error")
			}
			// gocritic: ifElseChain - rewrite if-else to switch statement
			switch {
			case tc.outputCheck != nil:
				tc.outputCheck(t, output)
			case tc.expect != "":
				assert.Equal(t, tc.expect, output, "Expected '%s'", tc.expect)
			default:
				assert.Equal(t, tc.expect, output, "Expected '%s'", tc.expect)
			}
		})
	}
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
		if err := f.Close(); err != nil {
			t.Fatalf("Failed to close test file: %v", err)
		}
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

func TestBashExecutor_RunCommand_AsyncBackground(t *testing.T) {
	executor := NewBashExecutor()
	// Use a command from AsyncCommands (e.g., 'open' or 'sleep' as a dummy)
	output, err := executor.RunCommand("open /tmp")
	// On macOS, 'open' should exist; on Linux/Windows, may not, so allow error but check output
	if err == nil {
		assert.Contains(t, output, "[launched open]", "Expected async launch message")
	}
	// Also test a non-async command to ensure it does not return async message
	output, err = executor.RunCommand("echo async-test")
	assert.NoError(t, err)
	assert.Equal(t, "async-test\n", output)
}

func TestIsAsyncCommand(t *testing.T) {
	tests := []struct {
		cmd      string
		expected bool
	}{
		{"idea .", true},
		{"code .", true},
		{"chrome", true},
		{"open /Applications/Calculator.app", true},
		{"echo hello", false},
		{"sleep 1", false},
		{"", false},
	}
	for _, tt := range tests {
		_, got := isAsyncCommand(tt.cmd)
		if got != tt.expected {
			t.Errorf("isAsyncCommand(%q) = %v, want %v", tt.cmd, got, tt.expected)
		}
	}
}

func TestBashExecutor_RunCommandWithDir_InteractiveAndAsync(t *testing.T) {
	executor := NewBashExecutor()
	dir := t.TempDir()

	// Async command (should return launch message)
	output, err := executor.RunCommandWithDir("open /tmp", dir)
	if err == nil {
		assert.Contains(t, output, "[launched open]")
	}

	// Interactive command (simulate with a non-blocking command)
	// We can't fully test PTY without a TTY, but we can check for error or empty output
	output, err = executor.RunCommandWithDir("vim --version", dir)
	// Accept error or empty output (since PTY may not work in test env)
	if err != nil {
		assert.Empty(t, output)
	}

	// Normal command
	output, err = executor.RunCommandWithDir("echo hi", dir)
	assert.NoError(t, err)
	assert.Equal(t, "hi\n", output)
}
