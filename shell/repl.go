package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// RunREPL starts an interactive read-eval-print loop
func RunREPL(sess *Session) error {
	scanner := bufio.NewScanner(os.Stdin)

	// Print the prompt before the first input
	fmt.Print(prompt(""))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			fmt.Print(prompt(""))
			continue
		}

		// Handle built-in exit commands
		if isExitCommand(line) {
			break
		}

		// Execute external command
		output, err := sess.Executor.RunCommand(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			// Print output without adding extra newlines
			fmt.Print(output)
		}

		// Only print prompt if output does not end with a newline
		if !strings.HasSuffix(output, "\n") {
			fmt.Print("\n")
		}
		fmt.Print(prompt(""))
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		// Don't treat EOF as an error - it's expected when user presses Ctrl-D
		if err == io.EOF {
			return nil
		}
		return err
	}

	return nil
}

// isExitCommand checks if the command is a built-in exit command
func isExitCommand(cmd string) bool {
	exitAliases := []string{"exit", "quit", "bye"}
	for _, alias := range exitAliases {
		if cmd == alias {
			return true
		}
	}
	return false
}
