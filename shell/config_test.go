package shell

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetColor(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"named color", "red", "\x1b[31m"},
		{"named color", "cyan", "\x1b[36m"},
		{"ansi code", "\x1b[35m", "\x1b[35m"},
		{"unknown", "notacolor", ""},
	}
	for _, c := range cases {
		if got := getColor(c.input); got != c.expected {
			t.Errorf("getColor(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestLoadColorConfig_EnvOverride(t *testing.T) {
	if err := os.Setenv("BINKS_PROMPT_COLOR", "yellow"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	if err := os.Setenv("BINKS_BRANCH_COLOR", "blue"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	if err := os.Setenv("BINKS_ERROR_COLOR", "green"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	cfg := LoadColorConfig()
	if cfg.PromptColor != "yellow" {
		t.Errorf("PromptColor = %q, want yellow", cfg.PromptColor)
	}
	if cfg.BranchColor != "blue" {
		t.Errorf("BranchColor = %q, want blue", cfg.BranchColor)
	}
	if cfg.ErrorColor != "green" {
		t.Errorf("ErrorColor = %q, want green", cfg.ErrorColor)
	}
	if err := os.Unsetenv("BINKS_PROMPT_COLOR"); err != nil {
		t.Fatalf("Unsetenv failed: %v", err)
	}
	if err := os.Unsetenv("BINKS_BRANCH_COLOR"); err != nil {
		t.Fatalf("Unsetenv failed: %v", err)
	}
	if err := os.Unsetenv("BINKS_ERROR_COLOR"); err != nil {
		t.Fatalf("Unsetenv failed: %v", err)
	}
}

func TestReadConfigFile_Errors(t *testing.T) {
	// Simulate error by temporarily renaming home dir
	home, _ := os.UserHomeDir()
	badPath := filepath.Join(home, ".binks.yaml")
	backup := badPath + ".bak"
	_ = os.Rename(badPath, backup) // ignore error if file doesn't exist
	cfg := readConfigFile()
	assert.Equal(t, ColorConfig{}, cfg)
	_ = os.Rename(backup, badPath) // restore

	// Simulate invalid YAML
	f, err := os.Create(badPath)
	assert.NoError(t, err)
	_, err = f.WriteString("not: [valid: yaml")
	assert.NoError(t, err)
	assert.NoError(t, f.Close())
	cfg = readConfigFile()
	assert.Equal(t, ColorConfig{}, cfg)
	_ = os.Remove(badPath)
}
