package shell

import (
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-isatty"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

const (
	PromptColor = "\x1b[36m" // Cyan
	ErrorColor  = "\x1b[31m" // Red
	ResetColor  = "\x1b[0m"
)

// StripANSI removes ANSI escape codes from a string (for test compatibility)
func StripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// formatPrompt returns a colored prompt string with current directory, using ~ for home.
func formatPrompt(cwd string) string {
	home, _ := os.UserHomeDir()
	shortCwd := cwd
	if home != "" && strings.HasPrefix(cwd, home) {
		shortCwd = "~" + cwd[len(home):]
	}
	return PromptColor + "binks:" + shortCwd + " > " + ResetColor + " "
}

// plainPrompt returns the prompt string without color codes
func plainPrompt(cwd string) string {
	home, _ := os.UserHomeDir()
	shortCwd := cwd
	if home != "" && strings.HasPrefix(cwd, home) {
		shortCwd = "~" + cwd[len(home):]
	}
	return "binks:" + shortCwd + " > "
}

// prompt returns the shell prompt string, colored if output is a TTY
func prompt(cwd string) string {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return formatPrompt(cwd)
	}
	return plainPrompt(cwd)
}
