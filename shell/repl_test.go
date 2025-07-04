package shell

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/binks-cli/binks/internal/executor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunREPL_Integration(t *testing.T) {
	// Test the REPL through the main binary as per acceptance criteria
	// This simulates: echo test\nexit\n input and checks output order

	binPath := "../binks"

	// Build the binary if it doesn't exist
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	require.NoError(t, buildCmd.Run(), "Failed to build binary")

	// Test the acceptance criteria: echo test\nexit\n
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader("echo test\nexit\n")
	output, err := cmd.CombinedOutput()

	require.NoError(t, err, "Expected no error from REPL")

	outputStr := string(output)

	// Check that we have prompts and output
	assert.Contains(t, outputStr, "binks>", "Expected prompt in output")
	// Should contain "test" output
	assert.Contains(t, outputStr, "test", "Expected 'test' output")
}

func TestIsExit(t *testing.T) {
	testCases := []struct {
		cmd      string
		expected bool
	}{
		// Basic cases
		{"exit", true},
		{"quit", true},
		{":q", true},

		// Case-insensitive matching
		{"EXIT", true},
		{"QUIT", true},
		{"Exit", true},
		{"Quit", true},
		{":Q", true},

		// With whitespace
		{"  exit  ", true},
		{"  quit  ", true},
		{"  :q  ", true},

		// Non-exit commands
		{"echo exit", false},
		{"", false},
		{"help", false},
		{"bye", false}, // removed from exit aliases
		{"exitnow", false},
		{"quitter", false},
		{":qa", false},
	}

	for _, tc := range testCases {
		result := isExit(tc.cmd)
		assert.Equal(t, tc.expected, result, "isExit(%q)", tc.cmd)
	}
}

func TestSession_NewSession(t *testing.T) {
	sess := NewSession()
	assert.NotNil(t, sess, "NewSession() returned nil")
	assert.NotNil(t, sess.Executor, "NewSession() created session with nil executor")
}

func TestRunREPL_MockExecutor(t *testing.T) {
	// Test with testify/mock executor for controlled testing
	mock := &executor.MockExecutorTestify{}
	mock.On("RunCommand", "echo hi").Return("hi\n", nil)
	mock.On("RunCommand", "failing-cmd").Return("", errors.New("command failed"))

	sess := &Session{Executor: mock}

	// Test that the session can use the mock executor
	output, err := sess.Executor.RunCommand("echo hi")
	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "hi", strings.TrimSpace(output), "Expected 'hi'")

	// Test error case
	_, err = sess.Executor.RunCommand("failing-cmd")
	assert.Error(t, err, "Expected error for failing command")

	mock.AssertExpectations(t)
}

func TestRunREPL_BlankLinePrompt(t *testing.T) {
	// Build the binary if it doesn't exist
	binPath := "../binks"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	_ = buildCmd.Run() // Ignore error if already built

	// Simulate pressing Enter (blank line), then exit
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader("\nexit\n")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Expected no error from REPL")
	outputStr := string(output)

	// Should see two prompts and nothing between them
	prompts := strings.Count(outputStr, "binks>")
	assert.GreaterOrEqual(t, prompts, 2, "Expected at least two prompts for blank line")

	// Should not see any output between the two prompts
	parts := strings.Split(outputStr, "binks>")
	if len(parts) >= 3 {
		assert.Equal(t, "", strings.TrimSpace(parts[1]), "Expected no output between prompts for blank line")
	}
}

func TestRunREPL_ExitHandling(t *testing.T) {
	// Test all exit commands with the main binary
	binPath := "../binks"

	// Build the binary if it doesn't exist
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	require.NoError(t, buildCmd.Run(), "Failed to build binary")

	testCases := []struct {
		name  string
		input string
	}{
		{"exit command", "exit\n"},
		{"quit command", "quit\n"},
		{"vim-style quit", ":q\n"},
		{"case insensitive EXIT", "EXIT\n"},
		{"case insensitive QUIT", "QUIT\n"},
		{"case insensitive :Q", ":Q\n"},
		{"exit with whitespace", "  exit  \n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binPath)
			cmd.Stdin = strings.NewReader(tc.input)
			output, err := cmd.CombinedOutput()

			require.NoError(t, err, "Expected clean exit for %s", tc.name)

			outputStr := string(output)

			// Should have at least one prompt
			assert.Contains(t, outputStr, "binks>", "Expected prompt in output for %s", tc.name)
		})
	}
}

func TestRunREPL_EOFHandling(t *testing.T) {
	// Test EOF handling (Ctrl-D)
	binPath := "../binks"

	// Build the binary if it doesn't exist
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	require.NoError(t, buildCmd.Run(), "Failed to build binary")

	// Create a command and close stdin immediately to simulate EOF
	cmd := exec.Command(binPath)

	// Create a pipe and close it immediately to simulate Ctrl-D (EOF)
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	// Start the command
	require.NoError(t, cmd.Start(), "Failed to start command")

	// Close stdin to send EOF
	stdin.Close()

	// Wait for the command to finish
	require.NoError(t, cmd.Wait(), "Expected clean exit on EOF")
}

func TestSession_ChangeDir(t *testing.T) {
	tmp := t.TempDir()
	sess := NewSession()

	// Change to temp dir
	err := sess.ChangeDir(tmp)
	require.NoError(t, err, "expected no error")
	cwdEval, _ := filepath.EvalSymlinks(sess.Cwd())
	tmpEval, _ := filepath.EvalSymlinks(tmp)
	assert.Equal(t, tmpEval, cwdEval, "expected cwd to match temp dir")

	// Change to parent dir
	err = sess.ChangeDir("..")
	require.NoError(t, err, "expected no error")
	assert.NotEqual(t, tmp, sess.Cwd(), "expected cwd to change from temp dir")

	// Change to home dir
	home, _ := os.UserHomeDir()
	homeEval, _ := filepath.EvalSymlinks(home)
	err = sess.ChangeDir("")
	require.NoError(t, err, "expected no error")
	cwdEval, _ = filepath.EvalSymlinks(sess.Cwd())
	assert.Equal(t, homeEval, cwdEval, "expected cwd to match home dir")

	// Invalid dir
	err = sess.ChangeDir("/no/such/dir/shouldexist")
	assert.Error(t, err, "expected error for invalid dir")
}
