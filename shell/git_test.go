package shell

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitBranch_NotInRepo(t *testing.T) {
	branch := GetGitBranch("/")
	if branch != "" {
		t.Errorf("Expected empty branch outside git repo, got %q", branch)
	}
}

func TestGetGitBranch_Detached(_ *testing.T) {
	// This test is a placeholder: to fully test, mock exec.Command or run in a temp git repo
	// For now, just ensure it doesn't panic
	_ = GetGitBranch("/")
}

func TestGetGitBranch_AndTrimNewline(t *testing.T) {
	// Not in a git repo: should return ""
	branch := GetGitBranch("/")
	assert.Equal(t, "", branch)

	// trimNewline removes trailing newlines
	assert.Equal(t, "foo", trimNewline("foo\n"))
	assert.Equal(t, "foo", trimNewline("foo\r\n"))
	assert.Equal(t, "foo", trimNewline("foo"))
	assert.Equal(t, "foo", trimNewline("foo\r")) // Now removes trailing \r
	assert.Equal(t, "", trimNewline("\n"))

	// Optionally: test in a temp git repo if git is available
	if _, err := exec.LookPath("git"); err == nil {
		dir := t.TempDir()
		cmd := exec.Command("git", "init")
		cmd.Dir = dir
		_ = cmd.Run()
		// Create a commit so branch exists
		_ = os.WriteFile(filepath.Join(dir, "README.md"), []byte("hi"), 0644)
		_ = exec.Command("git", "add", ".").Run()
		_ = exec.Command("git", "-C", dir, "add", ".").Run()
		_ = exec.Command("git", "-C", dir, "commit", "-m", "init").Run()
		branch := GetGitBranch(dir)
		assert.NotEmpty(t, branch)
	}
}

// More comprehensive tests would require a temp git repo and mocking exec.Command
