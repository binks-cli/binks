package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBashExecutor_RunCommand_SimpleEcho(t *testing.T) {
	executor := NewBashExecutor()
	
	output, err := executor.RunCommand("echo hello")
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if output != "hello\n" {
		t.Errorf("Expected 'hello\\n', got '%q'", output)
	}
}

func TestBashExecutor_RunCommand_NonExistentCommand(t *testing.T) {
	executor := NewBashExecutor()
	
	output, err := executor.RunCommand("nonexistentcommand12345")
	
	if err == nil {
		t.Fatal("Expected error for non-existent command, got none")
	}
	
	// Check that we got an error (exit status indicates command failure)
	if !strings.Contains(err.Error(), "exit status") {
		t.Errorf("Expected error message to contain 'exit status', got: %v", err)
	}
	
	// Output should contain the shell error message
	if !strings.Contains(output, "not found") && !strings.Contains(output, "command not found") {
		t.Errorf("Expected output to contain 'not found' or 'command not found', got: %s", output)
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
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		f.Close()
	}
	
	// List files in the temp directory
	output, err := executor.RunCommand("ls " + tempDir)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Check that all test files are in the output
	for _, file := range testFiles {
		if !strings.Contains(output, file) {
			t.Errorf("Expected output to contain '%s', got: %s", file, output)
		}
	}
	
	// Check that output contains multiple lines (newlines)
	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines of output, got %d: %s", len(lines), output)
	}
}

func TestBashExecutor_RunCommand_WithArguments(t *testing.T) {
	executor := NewBashExecutor()
	
	output, err := executor.RunCommand("echo 'hello world' | wc -w")
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Should return "2" with newline (word count)
	if strings.TrimSpace(output) != "2" {
		t.Errorf("Expected '2', got '%s'", strings.TrimSpace(output))
	}
}

func TestBashExecutor_RunCommand_EmptyCommand(t *testing.T) {
	executor := NewBashExecutor()
	
	output, err := executor.RunCommand("")
	
	// Empty command should succeed (bash -c "" returns 0)
	if err != nil {
		t.Fatalf("Expected no error for empty command, got: %v", err)
	}
	
	// Output should be empty
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}
}

func TestBashExecutor_RunCommand_ShellFeatures(t *testing.T) {
	executor := NewBashExecutor()
	
	// Test that shell features like wildcards work
	// Use a command that relies on shell expansion
	output, err := executor.RunCommand("echo $HOME")
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Should contain the home directory path
	if output == "$HOME" {
		t.Error("Expected shell variable expansion, but got literal '$HOME'")
	}
	
	// Should not be empty (unless HOME is somehow empty, which would be unusual)
	if output == "" {
		t.Error("Expected non-empty output for $HOME expansion")
	}
}
