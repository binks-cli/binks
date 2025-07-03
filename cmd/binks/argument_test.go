package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMainCLI_ArgumentHandling(t *testing.T) {
	// Test with arguments that contain spaces
	cmd := exec.Command("../../binks", "echo", "hello world")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if strings.TrimSpace(string(output)) != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", strings.TrimSpace(string(output)))
	}
}

func TestMainCLI_ArgumentsWithCommas(t *testing.T) {
	// Test with arguments that contain commas
	cmd := exec.Command("../../binks", "echo", "hello, world")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if strings.TrimSpace(string(output)) != "hello, world" {
		t.Errorf("Expected 'hello, world', got '%s'", strings.TrimSpace(string(output)))
	}
}

func TestMainCLI_ArgumentsWithSpecialChars(t *testing.T) {
	// Test with arguments that contain special characters (but not history expansion)
	cmd := exec.Command("../../binks", "echo", "hello@world#test")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if strings.TrimSpace(string(output)) != "hello@world#test" {
		t.Errorf("Expected 'hello@world#test', got '%s'", strings.TrimSpace(string(output)))
	}
}
