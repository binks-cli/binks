package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func containsPrompt(output string) bool {
	plain := "binks>"
	colored := "binks:"
	return strings.Contains(output, plain) || strings.Contains(output, colored)
}

func TestMainCLI_TableDriven(t *testing.T) {
	binPath := "../../binks"
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", binPath, ".")
		buildCmd.Dir = "../../"
		if err := buildCmd.Run(); err != nil {
			t.Fatalf("Failed to build binary: %v", err)
		}
	}

	testCases := []struct {
		name        string
		args        []string
		stdin       string
		expectError bool
		expect      string
	}{
		{
			name:   "echo command",
			args:   []string{"echo", "test"},
			stdin:  "",
			expect:  "test",
		},
		{
			name:   "REPL mode (no arguments)",
			args:   []string{},
			stdin:  "exit\n",
			expect:  "binks:",
		},
		{
			name:        "invalid command",
			args:        []string{"invalidcommand12345"},
			expectError: true,
			expect:      "Error:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tc.args...)
			if tc.stdin != "" {
				cmd.Stdin = strings.NewReader(tc.stdin)
			}
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tc.name == "REPL mode (no arguments)" {
				if !containsPrompt(outputStr) {
					t.Errorf("Expected prompt in output, got: %s", outputStr)
				}
				return
			}

			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if !strings.Contains(outputStr, tc.expect) {
					t.Errorf("Expected error message to contain '%s', got: %s", tc.expect, outputStr)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if !strings.Contains(outputStr, tc.expect) {
					t.Errorf("Expected output to contain '%s', got: %s", tc.expect, outputStr)
				}
			}
		})
	}
}
