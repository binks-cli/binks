package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockExecutor_RunCommand(t *testing.T) {
	mock := NewMockExecutor()

	// Test successful command
	mock.SetResponse("echo hello", "hello", nil)
	output, err := mock.RunCommand("echo hello")

	require.NoError(t, err, "Expected no error")

	assert.Equal(t, "hello", output, "Expected 'hello'")
}

func TestMockExecutor_RunCommand_WithError(t *testing.T) {
	mock := NewMockExecutor()

	// Test command that returns error
	expectedError := errors.New("command failed")
	mock.SetResponse("failing-command", "error output", expectedError)

	output, err := mock.RunCommand("failing-command")

	assert.Error(t, err, "Expected error")
	assert.Equal(t, expectedError, err, "Expected specific error")
	assert.Equal(t, "error output", output, "Expected 'error output'")
}

func TestMockExecutor_RunCommand_UnknownCommand(t *testing.T) {
	mock := NewMockExecutor()

	// Test command that wasn't set up
	output, err := mock.RunCommand("unknown-command")

	assert.Error(t, err, "Expected error for unknown command")
	assert.Equal(t, "", output, "Expected empty output")
}
