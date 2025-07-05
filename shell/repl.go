package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/binks-cli/binks/internal/agent"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
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
			Prompt:          promptWithAI(sess.Cwd(), sess.AIEnabled),
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
		return runREPLInteractive(sess, rl, os.Stdout, os.Stderr)
	}
	// Non-TTY: fallback to bufio.Scanner for integration tests and piping
	return RunREPLNonInteractive(sess, os.Stdin, os.Stdout, os.Stderr)
}

// RunREPLNonInteractive runs the REPL in non-interactive mode (for tests and piping).
// It reads lines from the provided io.Reader and writes output/errors to the given writers.
func RunREPLNonInteractive(sess *Session, in io.Reader, out, errOut io.Writer) error {
	scanner := bufio.NewScanner(in)
	// Print initial prompt
	fmt.Fprint(out, promptWithAI(sess.Cwd(), sess.AIEnabled))
	if f, ok := out.(interface{ Sync() error }); ok {
		_ = f.Sync()
	}
	for scanner.Scan() {
		line := scanner.Text()
		exit := processREPLLine(line, sess, out, errOut)
		// Print prompt after each command (to match interactive mode)
		fmt.Fprint(out, promptWithAI(sess.Cwd(), sess.AIEnabled))
		if f, ok := out.(interface{ Sync() error }); ok {
			_ = f.Sync()
		}
		if exit {
			break
		}
	}
	return scanner.Err()
}

// LineReader abstracts a line-oriented input for the REPL.
// It is implemented by *readline.Instance in production, and by mocks in tests.
type LineReader interface {
	// Readline reads the next line of input, or returns io.EOF at end.
	Readline() (string, error)
	// SetPrompt updates the prompt string.
	SetPrompt(string)
	// Close releases any resources held by the reader.
	Close() error
}

// runREPLInteractive runs the interactive REPL loop using a LineReader.
// This enables dependency injection for readline and testability.
func runREPLInteractive(sess *Session, rl LineReader, out, errOut io.Writer) error {
	defer rl.Close()
	for {
		line, err := rl.Readline()
		if err != nil {
			if err.Error() == "Interrupt" { // readline.ErrInterrupt is not exported
				if len(line) == 0 {
					break // exit on double Ctrl+C
				}
				continue
			} else if err == io.EOF {
				break // exit on Ctrl+D
			}
			return err
		}
		exit := processREPLLine(line, sess, out, errOut)
		if line == "" {
			rl.SetPrompt(promptWithAI(sess.Cwd(), sess.AIEnabled))
			continue
		}
		if strings.HasPrefix(line, "cd") || line == "help" || line == "?" {
			rl.SetPrompt(promptWithAI(sess.Cwd(), sess.AIEnabled))
			continue
		}
		rl.SetPrompt(promptWithAI(sess.Cwd(), sess.AIEnabled))
		if exit {
			break
		}
	}
	return nil
}

var aiColor = color.New(color.FgCyan, color.Bold)

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
	if strings.HasPrefix(line, ":ai ") {
		cmd := strings.TrimSpace(line[4:])
		if cmd == "on" {
			sess.AIEnabled = true
			fmt.Fprintln(out, "[AI mode enabled]")
		} else if cmd == "off" {
			sess.AIEnabled = false
			fmt.Fprintln(out, "[AI mode disabled]")
		} else {
			fmt.Fprintln(out, "Usage: :ai on|off")
		}
		return false
	}
	// AI query handling
	if sess.pendingSuggestion != nil {
		// We are waiting for user confirmation
		answer := strings.ToLower(strings.TrimSpace(line))
		if answer == "y" || answer == "yes" {
			sess.pendingSuggestion.confirmed = true
			output, err := sess.RunCommand(sess.pendingSuggestion.command)
			sess.pendingSuggestion = nil
			if err != nil {
				aiColor.Fprintf(out, "[AI] error: %s\n", err.Error())
			} else if output != "" {
				aiColor.Fprintf(out, "%s\n", output)
			}
		} else {
			sess.pendingSuggestion.declined = true
			aiColor.Fprintf(out, "[AI] Cancelled.\n")
			sess.pendingSuggestion = nil
		}
		return false
	}
	if sess.AIEnabled && sess.Agent != nil {
		if strings.HasPrefix(line, "!") {
			// Force shell command
			output, err := sess.RunCommand(strings.TrimSpace(line[1:]))
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
		resp, err := sess.ExecuteLine(agent.AIPrefix + line)
		if err != nil {
			aiColor.Fprintf(out, "[AI] error: %s\n", err.Error())
		} else if resp == "[AI]" && sess.pendingSuggestion != nil {
			// Show explanation and command, prompt for confirmation
			if sess.pendingSuggestion.explanation != "" {
				fmt.Fprintf(out, "[AI] %s\n", sess.pendingSuggestion.explanation)
			}
			fmt.Fprintf(out, "AI suggests: %s\n", sess.pendingSuggestion.command)
			fmt.Fprintf(out, "Execute this? [y/N]: ")
		} else {
			fmt.Fprintf(out, "%s\n", resp[5:])
		}
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

AI queries: Start your input with '>>' to ask the AI agent (e.g., '>> how do I list files?').
All other input is executed as shell commands in your shell environment.`
	if _, err := fmt.Fprintln(w, help); err != nil {
		fmt.Fprintln(os.Stderr, "failed to print help:", err)
	}
}

// promptWithAI returns the shell prompt string, with [AI] marker if AI mode is enabled.
func promptWithAI(cwd string, aiEnabled bool) string {
	if aiEnabled {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			return color.New(color.FgCyan, color.Bold).Sprintf("[AI] binks:%s > ", cwd)
		}
		return "[AI] binks:" + cwd + " > "
	}
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return formatPrompt(cwd)
	}
	return plainPrompt(cwd)
}
