package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

// BashExecutor implements the Executor interface using bash shell
type BashExecutor struct{}

// NewBashExecutor creates a new BashExecutor
func NewBashExecutor() *BashExecutor {
	return &BashExecutor{}
}

// isAsyncCommand returns true if the command should be run asynchronously (non-blocking)
func isAsyncCommand(cmd string) (string, bool) {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return "", false
	}
	for _, ac := range AsyncCommands {
		if fields[0] == ac {
			return ac, true
		}
	}
	return "", false
}

// RunCommandAsyncWithDir launches a command asynchronously (non-blocking)
func (e *BashExecutor) RunCommandAsyncWithDir(cmd string, dir string) (string, error) {
	execCmd := exec.Command("bash", "-c", cmd)
	if dir != "" {
		execCmd.Dir = dir
	}
	err := execCmd.Start()
	if err != nil {
		return "", err
	}
	// Optionally include PID: fmt.Sprintf("[launched %s (PID %d)]\n", ...)
	return fmt.Sprintf("[launched %s]\n", strings.Fields(cmd)[0]), nil
}

// RunCommandWithDir executes a command using bash in the specified directory and returns the combined output
func (e *BashExecutor) RunCommandWithDir(cmd string, dir string) (string, error) {
	if _, ok := isAsyncCommand(cmd); ok {
		return e.RunCommandAsyncWithDir(cmd, dir)
	}
	if isInteractiveCommand(cmd) {
		execCmd := exec.Command("bash", "-c", cmd)
		if dir != "" {
			execCmd.Dir = dir
		}
		ptmx, err := pty.Start(execCmd)
		if err != nil {
			return "", err
		}
		defer func() { _ = ptmx.Close() }()

		// Save the current terminal state
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

		// Handle terminal resize
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		go func() {
			for range ch {
				_ = pty.InheritSize(os.Stdin, ptmx)
			}
		}()
		ch <- syscall.SIGWINCH // Initial resize
		defer signal.Stop(ch)

		// Copy stdin/stdout
		go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
		_, _ = io.Copy(os.Stdout, ptmx)

		return "", execCmd.Wait()
	}
	execCmd := exec.Command("bash", "-c", cmd)
	if dir != "" {
		execCmd.Dir = dir
	}
	output, err := execCmd.CombinedOutput()
	return string(output), err // Preserve shell output including trailing newlines
}

// RunCommand executes a command using bash and returns the combined output
func (e *BashExecutor) RunCommand(cmd string) (string, error) {
	return e.RunCommandWithDir(cmd, "")
}
