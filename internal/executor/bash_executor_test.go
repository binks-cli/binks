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
		outputCheck func(t *testing.T, output string)
		errCheck    func(t *testing.T, err error)
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
			if tc.outputCheck != nil {
				tc.outputCheck(t, output)
			} else if tc.expect != "" {
				assert.Equal(t, tc.expect, output, "Expected '%s'", tc.expect)
			} else {
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
