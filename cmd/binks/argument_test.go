package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainCLI_ArgumentHandling(t *testing.T) {
	// Test with arguments that contain spaces
	cmd := exec.Command("../../binks", "echo", "hello world")
	output, err := cmd.CombinedOutput()

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "hello world", strings.TrimSpace(string(output)), "Expected 'hello world'")
}

func TestMainCLI_ArgumentsWithCommas(t *testing.T) {
	// Test with arguments that contain commas
	cmd := exec.Command("../../binks", "echo", "hello, world")
	output, err := cmd.CombinedOutput()

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "hello, world", strings.TrimSpace(string(output)), "Expected 'hello, world'")
}

func TestMainCLI_ArgumentsWithSpecialChars(t *testing.T) {
	// Test with arguments that contain special characters (but not history expansion)
	cmd := exec.Command("../../binks", "echo", "hello@world#test")
	output, err := cmd.CombinedOutput()

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "hello@world#test", strings.TrimSpace(string(output)), "Expected 'hello@world#test'")
}
