package executor

import (
	"errors"
	"testing"
)

func TestMockExecutor_RunCommand(t *testing.T) {
	mock := NewMockExecutor()
	
	// Test successful command
	mock.SetResponse("echo hello", "hello", nil)
	output, err := mock.RunCommand("echo hello")
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if output != "hello" {
		t.Errorf("Expected 'hello', got '%s'", output)
	}
}

func TestMockExecutor_RunCommand_WithError(t *testing.T) {
	mock := NewMockExecutor()
	
	// Test command that returns error
	expectedError := errors.New("command failed")
	mock.SetResponse("failing-command", "error output", expectedError)
	
	output, err := mock.RunCommand("failing-command")
	
	if err == nil {
		t.Error("Expected error, got none")
	}
	
	if err != expectedError {
		t.Errorf("Expected specific error, got: %v", err)
	}
	
	if output != "error output" {
		t.Errorf("Expected 'error output', got '%s'", output)
	}
}

func TestMockExecutor_RunCommand_UnknownCommand(t *testing.T) {
	mock := NewMockExecutor()
	
	// Test command that wasn't set up
	output, err := mock.RunCommand("unknown-command")
	
	if err == nil {
		t.Error("Expected error for unknown command, got none")
	}
	
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}
}
