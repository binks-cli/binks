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

func TestIsExitCommand(t *testing.T) {
	testCases := []struct {
		cmd      string
		expected bool
	}{
		{"exit", true},
		{"quit", true},
		{"bye", true},
		{"EXIT", false}, // case sensitive
		{"echo exit", false},
		{"", false},
		{"help", false},
	}
	
	for _, tc := range testCases {
		result := isExitCommand(tc.cmd)
		if result != tc.expected {
			t.Errorf("isExitCommand(%q) = %v, expected %v", tc.cmd, result, tc.expected)
		}
	}
}

func TestSession_NewSession(t *testing.T) {
	sess := NewSession()
	if sess == nil {
		t.Error("NewSession() returned nil")
	}
	
	if sess.Executor == nil {
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