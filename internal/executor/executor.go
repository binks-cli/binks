package executor

// Executor defines the interface for command execution
type Executor interface {
	RunCommand(cmd string) (string, error)
}
