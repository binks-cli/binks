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

// RunCommandWithDir executes a command using bash in the specified directory and returns the combined output
func (e *BashExecutor) RunCommandWithDir(cmd string, dir string) (string, error) {
	execCmd := exec.Command("bash", "-c", cmd)
	if dir != "" {
		execCmd.Dir = dir
	}
	output, err := execCmd.CombinedOutput()
	return string(output), err
}

// RunCommand executes a command using bash and returns the combined output
func (e *BashExecutor) RunCommand(cmd string) (string, error) {
	return e.RunCommandWithDir(cmd, "")
}
