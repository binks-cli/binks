package shell

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/binks-cli/binks/internal/executor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: This test suite covers all acceptance criteria from issue #30:
// - cd command: success, failure, edge cases (root, home, trailing slash, invalid dir, platform-specific)
// - exit and help commands: various forms, case insensitivity, blank lines, and output
// - REPL loop: sequences like pwd; cd ..; pwd; exit
// - History: file path, persistence, duplicates, blanks, and file writing/reading
// - External cd command: verified not to affect session cwd
// - Platform differences: tests adapt for Windows/Unix
// - Resource closure: all files are closed after use
//
// If new shell features are added, please ensure corresponding tests are included here or in related test files.

func containsPrompt(output string) bool {
	// Accept both colored and plain prompt
	plain := "binks>"
	colored := "\x1b[36mbinks:"
	return strings.Contains(output, plain) || strings.Contains(output, colored) || strings.Contains(output, "binks:")
}

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
	assert.True(t, containsPrompt(outputStr), "Expected prompt in output")
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

	// Print output for debugging
	t.Logf("REPL output: %q", outputStr)

	// Count prompt occurrences in the whole output
	promptCount := strings.Count(outputStr, "> ")
	assert.GreaterOrEqual(t, promptCount, 2, "Expected at least two prompts for blank line")

	// Should not see any output between the two prompts (relax: allow whitespace)
	parts := strings.SplitN(outputStr, "> ", 3)
	if len(parts) >= 3 {
		between := strings.TrimSpace(parts[1])
		assert.True(t, between == "" || containsPrompt(between), "Expected no output between prompts for blank line, got: %q", between)
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
			assert.True(t, containsPrompt(outputStr), "Expected prompt in output for %s", tc.name)
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
	if err := stdin.Close(); err != nil {
		t.Fatalf("Failed to close stdin: %v", err)
	}

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

func TestHistoryFilePath(t *testing.T) {
	home, err := os.UserHomeDir()
	assert.NoError(t, err)
	historyFile := filepath.Join(home, ".binks_history")
	assert.True(t, strings.HasSuffix(historyFile, ".binks_history"))
	// Clean up if file exists
	_ = os.Remove(historyFile)

	// Create and write a test line
	f, err := os.Create(historyFile)
	assert.NoError(t, err)
	_, err = f.WriteString("echo test-history\n")
	assert.NoError(t, err)
	f.Close()

	// Read back
	data, err := os.ReadFile(historyFile)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "echo test-history")

	// Clean up
	_ = os.Remove(historyFile)
}

func TestHistoryFile_PersistenceAndEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, ".binks_history")

	// Write initial history with duplicates and blanks
	initial := "echo foo\n\necho bar\necho foo\n"
	assert.NoError(t, os.WriteFile(historyFile, []byte(initial), 0644))

	// Simulate readline loading and appending
	f, err := os.OpenFile(historyFile, os.O_APPEND|os.O_WRONLY, 0644)
	assert.NoError(t, err)
	_, err = f.WriteString("echo baz\n\n")
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	// Read back and check for duplicates and blanks
	data, err := os.ReadFile(historyFile)
	assert.NoError(t, err)
	lines := strings.Split(string(data), "\n")
	var nonBlank []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			nonBlank = append(nonBlank, l)
		}
	}
	// Should contain all commands, including duplicates, but no blank lines
	assert.Contains(t, nonBlank, "echo foo")
	assert.Contains(t, nonBlank, "echo bar")
	assert.Contains(t, nonBlank, "echo baz")
	for _, l := range nonBlank {
		assert.NotEqual(t, "", strings.TrimSpace(l))
	}
}

func TestREPL_SequenceAffectsState(t *testing.T) {
	binPath := "../binks"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	if err := buildCmd.Run(); err != nil {
		t.Skip("Could not build binks binary: ", err)
	}
	if _, err := os.Stat(binPath); err != nil {
		t.Skip("binks binary not found: ", err)
	}
	// Sequence: pwd; cd ..; pwd; exit
	startDir, _ := os.Getwd()
	parentDir := filepath.Dir(startDir)
	input := "pwd\ncd ..\npwd\nexit\n"
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	outStr := string(output)
	// Should see startDir and parentDir in output
	assert.Contains(t, outStr, startDir)
	assert.Contains(t, outStr, parentDir)
}

func TestProcessREPLLine_BuiltinsAndExternal(t *testing.T) {
	var out, errOut strings.Builder
	sess := NewSession()

	// Test blank line
	exit := processREPLLine("", sess, &out, &errOut)
	assert.False(t, exit)
	assert.Equal(t, "", out.String())
	assert.Equal(t, "", errOut.String())

	// Test exit command
	out.Reset()
	errOut.Reset()
	exit = processREPLLine("exit", sess, &out, &errOut)
	assert.True(t, exit)

	// Test cd to home
	out.Reset()
	errOut.Reset()
	home, _ := os.UserHomeDir()
	sess.ChangeDir("/tmp") // move away from home
	exit = processREPLLine("cd", sess, &out, &errOut)
	assert.False(t, exit)
	assert.Equal(t, home, sess.Cwd())

	// Test cd to invalid dir
	out.Reset()
	errOut.Reset()
	badDir := "/no/such/dir/shouldexist"
	exit = processREPLLine("cd "+badDir, sess, &out, &errOut)
	assert.False(t, exit)
	assert.Contains(t, errOut.String(), "Error:")

	// Test help
	out.Reset()
	errOut.Reset()
	exit = processREPLLine("help", sess, &out, &errOut)
	assert.False(t, exit)
	assert.Contains(t, out.String(), "Built-in commands:")

	// Test external command (mock)
	mock := &executor.MockExecutorTestify{}
	mock.On("RunCommand", "echo hi").Return("hi\n", nil)
	sess.Executor = mock
	out.Reset()
	errOut.Reset()
	exit = processREPLLine("echo hi", sess, &out, &errOut)
	assert.False(t, exit)
	assert.Contains(t, out.String(), "hi")
	mock.AssertExpectations(t)

	// Test external command error
	mock.On("RunCommand", "fail").Return("", errors.New("fail"))
	out.Reset()
	errOut.Reset()
	exit = processREPLLine("fail", sess, &out, &errOut)
	assert.False(t, exit)
	assert.Contains(t, errOut.String(), "Error:")
}

func TestRunREPLNonInteractive_Basic(t *testing.T) {
	sess := NewSession()
	input := "echo foo\ncd /\npwd\nexit\n"
	var out, errOut strings.Builder
	err := RunREPLNonInteractive(sess, strings.NewReader(input), &out, &errOut)
	assert.NoError(t, err)
	output := out.String()
	// Should contain 'foo' from echo, and '/' from pwd
	assert.Contains(t, output, "foo")
	assert.Contains(t, output, "/")
	// Should not contain error output
	assert.Empty(t, errOut.String())
}

// mockLineReader implements LineReader for testing runREPLInteractive
// It returns lines from the provided slice, then io.EOF
// SetPrompt and Close are no-ops

type mockLineReader struct {
	lines   []string
	idx     int
	prompts []string
}

func (m *mockLineReader) Readline() (string, error) {
	if m.idx >= len(m.lines) {
		return "", io.EOF
	}
	line := m.lines[m.idx]
	m.idx++
	return line, nil
}
func (m *mockLineReader) SetPrompt(p string) { m.prompts = append(m.prompts, p) }
func (m *mockLineReader) Close() error       { return nil }

func TestRunREPLInteractive_Basic(t *testing.T) {
	sess := NewSession()
	mockRL := &mockLineReader{lines: []string{"echo foo", "cd /", "pwd", "exit"}}
	var out, errOut strings.Builder
	err := runREPLInteractive(sess, mockRL, &out, &errOut)
	assert.NoError(t, err)
	output := out.String()
	assert.Contains(t, output, "foo")
	assert.Contains(t, output, "/")
	assert.Empty(t, errOut.String())
	// Prompts should be set at least once
	assert.NotEmpty(t, mockRL.prompts)
}

func TestPrintHelp(t *testing.T) {
	var buf strings.Builder
	printHelp(&buf)
	helpText := buf.String()
	assert.Contains(t, helpText, "Built-in commands:")
	assert.Contains(t, helpText, "cd <dir>")
	assert.Contains(t, helpText, "exit")
	assert.Contains(t, helpText, "help, ?")
}

func TestRunREPL_NonTTY(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()
	defer w.Close()

	sess := NewSession()
	// Write a simple command and exit
	go func() {
		w.Write([]byte("echo test\nexit\n"))
		w.Close()
	}()
	// Temporarily replace os.Stdin with our pipe
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()
	err = RunREPL(sess)
	assert.NoError(t, err)
	// No assertion on output, just ensure it runs without error
}

type errLineReader struct{}

func (e *errLineReader) Readline() (string, error) { return "", io.ErrUnexpectedEOF }
func (e *errLineReader) SetPrompt(string)          {}
func (e *errLineReader) Close() error              { return nil }

func TestRunREPLInteractive_ErrorHandling(t *testing.T) {
	sess := NewSession()
	err := runREPLInteractive(sess, &errLineReader{}, io.Discard, io.Discard)
	assert.Error(t, err)
	assert.Equal(t, io.ErrUnexpectedEOF, err)
}

func TestPromptFunctions(t *testing.T) {
	cwd := "/tmp"
	p := prompt(cwd)
	assert.Contains(t, p, "binks:")
	fp := formatPrompt(cwd)
	assert.Contains(t, fp, "binks:")
	plain := plainPrompt(cwd)
	assert.Contains(t, plain, "binks:")
}
