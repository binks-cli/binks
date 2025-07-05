package executor

import "testing"

func TestIsInteractiveCommand(t *testing.T) {
	cases := []struct {
		cmd      string
		expected bool
	}{
		{"vim", true},
		{"vim file.txt", true},
		{"nano foo", true},
		{"less bar", true},
		{"ssh user@host", true},
		{"ls -l", false},
		{"echo hello", false},
		{"cat file", false},
		{"man ls", true},
		{"top", true},
		{"htop", true},
		{"nvim", true},
		{"vi", true},
	}
	for _, c := range cases {
		if got := isInteractiveCommand(c.cmd); got != c.expected {
			t.Errorf("isInteractiveCommand(%q) = %v, want %v", c.cmd, got, c.expected)
		}
	}
}
