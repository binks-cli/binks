package shell

import "github.com/binks-cli/binks/internal/executor"

// Session represents the state of a shell session
type Session struct {
	Executor executor.Executor
	// Future fields for working directory, history, etc.
}

// NewSession creates a new shell session
func NewSession() *Session {
	return &Session{
		Executor: executor.NewBashExecutor(),
	}
}