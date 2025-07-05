package shell

import (
	"os"
	"strings"
	"testing"
)

func TestFormatPrompt(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		cwd      string
		expected string
	}{
		{cwd: home, expected: "~"},
		{cwd: home + "/project", expected: "~/project"},
		{cwd: "/tmp", expected: "/tmp"},
	}
	for _, tt := range tests {
		prompt := formatPrompt(tt.cwd)
		if !strings.Contains(prompt, tt.expected) {
			t.Errorf("prompt %q does not contain expected path %q", prompt, tt.expected)
		}
		if !strings.HasPrefix(prompt, "\x1b[") {
			t.Errorf("prompt %q does not start with ANSI code", prompt)
		}
		if !strings.HasSuffix(prompt, "\x1b[0m ") {
			t.Errorf("prompt %q does not end with reset code and space", prompt)
		}
	}
}
