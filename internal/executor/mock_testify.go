package executor

import "github.com/stretchr/testify/mock"

// MockExecutorTestify is a testify-based mock for the Executor interface
// Used to replace manual MockExecutor in tests
type MockExecutorTestify struct {
	mock.Mock
}

// RunCommand mocks the Executor's RunCommand method.
func (m *MockExecutorTestify) RunCommand(cmd string) (string, error) {
	args := m.Called(cmd)
	return args.String(0), args.Error(1)
}
