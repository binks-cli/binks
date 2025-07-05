// Package shell provides shell utilities and configuration for binks.
package shell

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ColorConfig holds color settings for the prompt, branch, and error messages.
type ColorConfig struct {
	PromptColor string `yaml:"prompt_color"`
	BranchColor string `yaml:"branch_color"`
	ErrorColor  string `yaml:"error_color"`
	// Future: add MCP, editor, etc.
}

// BinksConfig holds the overall configuration for binks.
type BinksConfig struct {
	Colors ColorConfig `yaml:"colors"`
	// Future: MCP, editor, etc.
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

// LoadColorConfig loads color config from YAML file and env vars, falling back to defaults
func LoadColorConfig() ColorConfig {
	cfg := defaultColors
	fileCfg := readConfigFile()
	if fileCfg.PromptColor != "" {
		cfg.PromptColor = fileCfg.PromptColor
	}
	if fileCfg.BranchColor != "" {
		cfg.BranchColor = fileCfg.BranchColor
	}
	if fileCfg.ErrorColor != "" {
		cfg.ErrorColor = fileCfg.ErrorColor
	}
	if v := os.Getenv("BINKS_PROMPT_COLOR"); v != "" {
		cfg.PromptColor = v
	}
	if v := os.Getenv("BINKS_BRANCH_COLOR"); v != "" {
		cfg.BranchColor = v
	}
	if v := os.Getenv("BINKS_ERROR_COLOR"); v != "" {
		cfg.ErrorColor = v
	}
	return cfg
}

// readConfigFile loads ~/.binks.yaml if present
func readConfigFile() ColorConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		return ColorConfig{}
	}
	path := filepath.Join(home, ".binks.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return ColorConfig{}
	}
	var cfg BinksConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return ColorConfig{}
	}
	return cfg.Colors
}
