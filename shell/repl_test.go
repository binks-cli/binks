package shell

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/binks-cli/binks/internal/executor"
)

func TestRunREPL_Integration(t *testing.T) {
	// Test the REPL through the main binary as per acceptance criteria
	// This simulates: echo test\nexit\n input and checks output order

	binPath := "../binks"

	// Build the binary if it doesn't exist
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Test the acceptance criteria: echo test\nexit\n
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader("echo test\nexit\n")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Expected no error from REPL, got: %v", err)
	}

	outputStr := string(output)

	// Check that we have prompts and output
	if !strings.Contains(outputStr, "binks>") {
		t.Errorf("Expected prompt in output, got: %s", outputStr)
	}

	// Should contain "test" output
	if !strings.Contains(outputStr, "test") {
		t.Errorf("Expected 'test' output, got: %s", outputStr)
	}
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
		if result != tc.expected {
			t.Errorf("isExit(%q) = %v, expected %v", tc.cmd, result, tc.expected)
		}
	}
}

func TestSession_NewSession(t *testing.T) {
	sess := NewSession()
	if sess == nil {
		t.Error("NewSession() returned nil")
	} else if sess.Executor == nil {
		t.Error("NewSession() created session with nil executor")
	}
}

func TestRunREPL_MockExecutor(t *testing.T) {
	// Test with mock executor for controlled testing
	mock := executor.NewMockExecutor()
	mock.SetResponse("echo hi", "hi\n", nil)
	mock.SetResponse("failing-cmd", "", errors.New("command failed"))

	sess := &Session{Executor: mock}

	// Test that the session can use the mock executor
	output, err := sess.Executor.RunCommand("echo hi")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if strings.TrimSpace(output) != "hi" {
		t.Errorf("Expected 'hi', got '%s'", strings.TrimSpace(output))
	}

	// Test error case
	_, err = sess.Executor.RunCommand("failing-cmd")
	if err == nil {
		t.Error("Expected error for failing command")
	}
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
	if err != nil {
		t.Fatalf("Expected no error from REPL, got: %v", err)
	}
	outputStr := string(output)

	// Should see two prompts and nothing between them
	prompts := strings.Count(outputStr, "binks>")
	if prompts < 2 {
		t.Errorf("Expected at least two prompts for blank line, got: %s", outputStr)
	}

	// Should not see any output between the two prompts
	parts := strings.Split(outputStr, "binks>")
	if len(parts) >= 3 && strings.TrimSpace(parts[1]) != "" {
		t.Errorf("Expected no output between prompts for blank line, got: %q", parts[1])
	}
}

func TestRunREPL_ExitHandling(t *testing.T) {
	// Test all exit commands with the main binary
	binPath := "../binks"

	// Build the binary if it doesn't exist
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

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

			// Process should exit cleanly with status 0
			if err != nil {
				t.Fatalf("Expected clean exit for %s, got error: %v", tc.name, err)
			}

			outputStr := string(output)

			// Should have at least one prompt
			if !strings.Contains(outputStr, "binks>") {
				t.Errorf("Expected prompt in output for %s, got: %s", tc.name, outputStr)
			}
		})
	}
}

func TestRunREPL_EOFHandling(t *testing.T) {
	// Test EOF handling (Ctrl-D)
	binPath := "../binks"

	// Build the binary if it doesn't exist
	buildCmd := exec.Command("go", "build", "-o", binPath, "../cmd/binks")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Create a command and close stdin immediately to simulate EOF
	cmd := exec.Command(binPath)

	// Create a pipe and close it immediately to simulate Ctrl-D (EOF)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// Close stdin to send EOF
	stdin.Close()

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		t.Fatalf("Expected clean exit on EOF, got error: %v", err)
	}
}
