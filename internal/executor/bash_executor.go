package executor

import (
	"os/exec"
)

// BashExecutor implements the Executor interface using bash shell
type BashExecutor struct{}

// NewBashExecutor creates a new BashExecutor
func NewBashExecutor() *BashExecutor {
	return &BashExecutor{}
}

// RunCommand executes a command using bash and returns the combined output
func (e *BashExecutor) RunCommand(cmd string) (string, error) {
	// Use bash -c to execute the command, enabling shell features like globs, aliases, etc.
	execCmd := exec.Command("bash", "-c", cmd)
	
	// Capture both stdout and stderr
	output, err := execCmd.CombinedOutput()
	
	// Return output as string, preserving the natural newlines for proper terminal display
	return string(output), err
}
