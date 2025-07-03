package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainCLI_EchoCommand(t *testing.T) {
	// Build the binary if it doesn't exist
	binPath := "../../binks"
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", binPath, ".")
		buildCmd.Dir = "../../"
		if err := buildCmd.Run(); err != nil {
			t.Fatalf("Failed to build binary: %v", err)
		}
	}

	// Test echo command
	cmd := exec.Command(binPath, "echo", "test")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if strings.TrimSpace(string(output)) != "test" {
		t.Errorf("Expected 'test', got '%s'", strings.TrimSpace(string(output)))
	}
}

func TestMainCLI_NoArguments(t *testing.T) {
	// Test with no arguments - should start REPL mode
	binPath := "../../binks"
	cmd := exec.Command(binPath)
	
	// Provide input to exit the REPL immediately
	cmd.Stdin = strings.NewReader("exit\n")
	output, err := cmd.CombinedOutput()
	
	// REPL should exit cleanly when receiving "exit" command
	if err != nil {
		t.Errorf("Expected no error for REPL mode, got: %v", err)
	}
	
	outputStr := string(output)
	if !strings.Contains(outputStr, "binks>") {
		t.Errorf("Expected REPL prompt, got: %s", outputStr)
	}
}

func TestMainCLI_InvalidCommand(t *testing.T) {
	// Test with invalid command
	binPath := "../../binks"
	cmd := exec.Command(binPath, "invalidcommand12345")
	output, err := cmd.CombinedOutput()
	
	// Should exit with error
	if err == nil {
		t.Error("Expected error for invalid command")
	}
	
	outputStr := string(output)
	if !strings.Contains(outputStr, "Error:") {
		t.Errorf("Expected error message, got: %s", outputStr)
	}
}
