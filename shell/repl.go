package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/mattn/go-isatty"
)

// RunREPL starts an interactive read-eval-print loop
func RunREPL(sess *Session) error {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		// Use readline for interactive TTY
		config := &readline.Config{
			Prompt:          prompt(sess.Cwd()),
			HistoryLimit:    100,
			InterruptPrompt: "^C\n",
			EOFPrompt:       "exit\n",
			Stdin:           os.Stdin,
			Stdout:          os.Stdout,
		}
		rl, err := readline.NewEx(config)
		if err != nil {
			return err
		}
		defer rl.Close()
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
			line = strings.TrimSpace(line)
			if line == "" {
				rl.SetPrompt(prompt(sess.Cwd()))
				continue
			}
			if isExit(line) {
				break
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
					fmt.Fprint(os.Stderr, ErrorMessage(err))
				}
				rl.SetPrompt(formatPrompt(sess.Cwd()))
				continue
			}
			if line == "help" || line == "?" {
				printHelp(os.Stdout)
				rl.SetPrompt(formatPrompt(sess.Cwd()))
				continue
			}
			output, err := sess.RunCommand(line)
			if err != nil {
				fmt.Fprint(os.Stderr, ErrorMessage(err))
			} else if output != "" {
				fmt.Print(output)
				if !strings.HasSuffix(output, "\n") {
					fmt.Print("\n")
				}
			}
			rl.SetPrompt(formatPrompt(sess.Cwd()))
		}
		return nil
	}
	// Non-TTY: fallback to bufio.Scanner for integration tests and piping
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt(sess.Cwd()))
	os.Stdout.Sync()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Print(prompt(sess.Cwd()))
			os.Stdout.Sync()
			continue
		}
		if isExit(line) {
			break
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
				fmt.Fprint(os.Stderr, ErrorMessage(err))
			}
			fmt.Print(formatPrompt(sess.Cwd()))
			os.Stdout.Sync()
			continue
		}
		if line == "help" || line == "?" {
			printHelp(os.Stdout)
			fmt.Print(formatPrompt(sess.Cwd()))
			os.Stdout.Sync()
			continue
		}
		output, err := sess.RunCommand(line)
		if err != nil {
			fmt.Fprint(os.Stderr, ErrorMessage(err))
		} else if output != "" {
			fmt.Print(output)
			if !strings.HasSuffix(output, "\n") {
				fmt.Print("\n")
			}
		}
		fmt.Print(formatPrompt(sess.Cwd()))
		os.Stdout.Sync()
	}
	if err := scanner.Err(); err != nil {
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

// printHelp prints the built-in help message to the given writer
func printHelp(w io.Writer) {
	help := `Built-in commands:
  cd <dir>    – Change directory
  exit        – Exit the shell
  help, ?     – Show this help message

All other input is executed as shell commands in your shell environment.`
	fmt.Fprintln(w, help)
}
