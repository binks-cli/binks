package shell

import (
	"os"
	"strings"
)

type ColorConfig struct {
	PromptColor string
	BranchColor string
	ErrorColor  string
}

var defaultColors = ColorConfig{
	PromptColor: "cyan",
	BranchColor: "magenta",
	ErrorColor:  "red",
}

var colorMap = map[string]string{
	"black":   "\x1b[30m",
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"magenta": "\x1b[35m",
	"cyan":    "\x1b[36m",
	"white":   "\x1b[37m",
	"reset":   "\x1b[0m",
}

// getColor returns the ANSI code for a color name or code
func getColor(name string) string {
	if code, ok := colorMap[strings.ToLower(name)]; ok {
		return code
	}
	if strings.HasPrefix(name, "\x1b[") {
		return name // already an ANSI code
	}
	return "" // fallback: no color
}

// LoadColorConfig loads color config from env vars, falling back to defaults
func LoadColorConfig() ColorConfig {
	cfg := defaultColors
	if v := os.Getenv("BINKS_PROMPT_COLOR"); v != "" {
		cfg.PromptColor = v
	}
	if v := os.Getenv("BINKS_BRANCH_COLOR"); v != "" {
		cfg.BranchColor = v
	}
	if v := os.Getenv("BINKS_ERROR_COLOR"); v != "" {
		cfg.ErrorColor = v
	}
	// TODO: add YAML config file support
	return cfg
}
