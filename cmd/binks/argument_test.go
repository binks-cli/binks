package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainCLI_ArgumentHandling_TableDriven(t *testing.T) {
	testCases := []struct {
		name     string
		arg      string
		expected string
	}{
		{"spaces", "hello world", "hello world"},
		{"commas", "hello, world", "hello, world"},
		{"special chars", "hello@world#test", "hello@world#test"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("../../binks", "echo", tc.arg)
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Expected no error")
			assert.Equal(t, tc.expected, strings.TrimSpace(string(output)), "Expected '%s'", tc.expected)
		})
	}
}
