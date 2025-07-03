package executor

import (
	"fmt"
	"strings"
)

// MockExecutor is a test implementation of the Executor interface
type MockExecutor struct {
	responses map[string]MockResponse
}

// MockResponse represents a mocked command response
type MockResponse struct {
	Output string
	Error  error
}

// NewMockExecutor creates a new MockExecutor
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		responses: make(map[string]MockResponse),
	}
}

// SetResponse sets the expected response for a given command
func (m *MockExecutor) SetResponse(cmd string, output string, err error) {
	m.responses[cmd] = MockResponse{Output: output, Error: err}
}

// RunCommand executes a mocked command
func (m *MockExecutor) RunCommand(cmd string) (string, error) {
	response, exists := m.responses[cmd]
	if !exists {
		return "", fmt.Errorf("mock: command not found: %s", cmd)
	}
	
	return response.Output, response.Error
}

// SetDefaultEchoResponse sets up the mock to handle echo commands
func (m *MockExecutor) SetDefaultEchoResponse() {
	// Handle echo commands by returning the text after "echo "
	for cmd := range m.responses {
		if strings.HasPrefix(cmd, "echo ") {
			text := strings.TrimPrefix(cmd, "echo ")
			m.responses[cmd] = MockResponse{Output: text, Error: nil}
		}
	}
}
