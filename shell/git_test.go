package shell

import (
	"testing"
)

func TestGetGitBranch_NotInRepo(t *testing.T) {
	branch := GetGitBranch("/")
	if branch != "" {
		t.Errorf("Expected empty branch outside git repo, got %q", branch)
	}
}

func TestGetGitBranch_Detached(t *testing.T) {
	// This test is a placeholder: to fully test, mock exec.Command or run in a temp git repo
	// For now, just ensure it doesn't panic
	_ = GetGitBranch("/")
}

// More comprehensive tests would require a temp git repo and mocking exec.Command
