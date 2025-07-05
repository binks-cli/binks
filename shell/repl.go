package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/mattn/go-isatty"
)

// RunREPL starts an interactive read-eval-print loop
func RunREPL(sess *Session) error {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		// Use readline for interactive TTY
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		historyFile := filepath.Join(homeDir, ".binks_history")
		config := &readline.Config{
			Prompt:          prompt(sess.Cwd()),
			HistoryLimit:    100,
			InterruptPrompt: "^C\n",
			EOFPrompt:       "exit\n",
			Stdin:           os.Stdin,
			Stdout:          os.Stdout,
			HistoryFile:     historyFile,
		}
		rl, err := readline.NewEx(config)
		if err != nil {
			return err
		}
		defer func() {
			if err := rl.Close(); err != nil {
				fmt.Fprintln(os.Stderr, "failed to close readline:", err)
			}
		}()
		for {
			line, err := rl.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break // exit on double Ctrl+C
				}
				continue
			} else if err == io.EOF {
				break // exit on Ctrl+D
			}
			exit := processREPLLine(line, sess, os.Stdout, os.Stderr)
			if line == "" {
				rl.SetPrompt(prompt(sess.Cwd()))
				continue
			}
			if strings.HasPrefix(line, "cd") || line == "help" || line == "?" {
				rl.SetPrompt(formatPrompt(sess.Cwd()))
				continue
			}
			rl.SetPrompt(formatPrompt(sess.Cwd()))
			if exit {
				break
			}
		}
		return nil
	}
	// Non-TTY: fallback to bufio.Scanner for integration tests and piping
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt(sess.Cwd()))
	if err := os.Stdout.Sync(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to sync stdout:", err)
	}
	for scanner.Scan() {
		line := scanner.Text()
		exit := processREPLLine(line, sess, os.Stdout, os.Stderr)
		if strings.TrimSpace(line) == "" {
			fmt.Print(prompt(sess.Cwd()))
			if err := os.Stdout.Sync(); err != nil {
				fmt.Fprintln(os.Stderr, "failed to sync stdout:", err)
			}
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "cd") || strings.TrimSpace(line) == "help" || strings.TrimSpace(line) == "?" {
			fmt.Print(formatPrompt(sess.Cwd()))
			if err := os.Stdout.Sync(); err != nil {
				fmt.Fprintln(os.Stderr, "failed to sync stdout:", err)
			}
			continue
		}
		fmt.Print(formatPrompt(sess.Cwd()))
		if err := os.Stdout.Sync(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to sync stdout:", err)
		}
		if exit {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}

// processREPLLine handles a single REPL input line and returns whether to exit the loop.
func processREPLLine(line string, sess *Session, out, errOut io.Writer) (exit bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}
	if isExit(line) {
		return true
	}
	if strings.HasPrefix(line, "cd") {
		fields := strings.Fields(line)
		var cdArg string
		if len(fields) > 1 {
			cdArg = strings.Join(fields[1:], " ")
		} else {
			cdArg = ""
		}
		err := sess.ChangeDir(strings.TrimSpace(cdArg))
		if err != nil {
			fmt.Fprint(errOut, ErrorMessage(err))
		}
		return false
	}
	if line == "help" || line == "?" {
		printHelp(out)
		return false
	}
	output, err := sess.RunCommand(line)
	if err != nil {
		fmt.Fprint(errOut, ErrorMessage(err))
	} else if output != "" {
		fmt.Fprint(out, output)
		if !strings.HasSuffix(output, "\n") {
			fmt.Fprint(out, "\n")
		}
	}
	return false
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

// printHelp prints the built-in help message to the given writer
func printHelp(w io.Writer) {
	help := `Built-in commands:
  cd <dir>    – Change directory
  exit        – Exit the shell
  help, ?     – Show this help message

All other input is executed as shell commands in your shell environment.`
	if _, err := fmt.Fprintln(w, help); err != nil {
		fmt.Fprintln(os.Stderr, "failed to print help:", err)
	}
}
