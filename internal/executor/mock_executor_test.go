package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockExecutorTestify_RunCommand(t *testing.T) {
	mock := &MockExecutorTestify{}
	mock.On("RunCommand", "echo hello").Return("hello", nil)

	output, err := mock.RunCommand("echo hello")
	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "hello", output, "Expected 'hello'")
	mock.AssertExpectations(t)
}

func TestMockExecutorTestify_RunCommand_WithError(t *testing.T) {
	mock := &MockExecutorTestify{}
	expectedError := errors.New("command failed")
	mock.On("RunCommand", "failing-command").Return("error output", expectedError)

	output, err := mock.RunCommand("failing-command")
	assert.Error(t, err, "Expected error")
	assert.Equal(t, expectedError, err, "Expected specific error")
	assert.Equal(t, "error output", output, "Expected 'error output'")
	mock.AssertExpectations(t)
}

func TestMockExecutorTestify_RunCommand_UnknownCommand(t *testing.T) {
	mock := &MockExecutorTestify{}
	mock.On("RunCommand", "unknown-command").Return("", errors.New("mock: command not found: unknown-command"))

	output, err := mock.RunCommand("unknown-command")
	assert.Error(t, err, "Expected error for unknown command")
	assert.Equal(t, "", output, "Expected empty output")
	mock.AssertExpectations(t)
}
