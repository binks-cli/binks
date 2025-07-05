package shell

import (
	"os/exec"
)

// GetGitBranch returns the current git branch name, or short commit hash if detached, or empty string if not in a git repo.
func GetGitBranch(cwd string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	branch := string(out)
	branch = trimNewline(branch)
	if branch == "HEAD" {
		// Detached HEAD, get short commit hash
		cmd2 := exec.Command("git", "rev-parse", "--short", "HEAD")
		cmd2.Dir = cwd
		out2, err2 := cmd2.Output()
		if err2 != nil {
			return "detached"
		}
		return trimNewline(string(out2))
	}
	return branch
}

func trimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}
