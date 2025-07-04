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
		if isExit(line) {
			break
		}

		// Execute external command
		output, err := sess.Executor.RunCommand(line)
		if err != nil {
			errMsg := fmt.Sprintf("Error: %v", err)
			fmt.Fprint(os.Stderr, errMsg)
			if !strings.HasSuffix(errMsg, "\n") {
				fmt.Fprint(os.Stderr, "\n")
			}
		} else if output != "" {
			// Print output only if there is output
			fmt.Print(output)
			// Add newline only if output doesn't end with one
			if !strings.HasSuffix(output, "\n") {
				fmt.Print("\n")
			}
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

// isExit checks if the command is a built-in exit command (case-insensitive)
func isExit(line string) bool {
	// Convert to lowercase for case-insensitive matching
	cmd := strings.ToLower(strings.TrimSpace(line))
	exitAliases := []string{"exit", "quit", ":q"}
	for _, alias := range exitAliases {
		if cmd == alias {
			return true
		}
	}
	return false
}
