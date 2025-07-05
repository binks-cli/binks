package agent

import "strings"

const aiPrefix = ">>"

// IsAIQuery returns true if the line is an AI query (starts with the AI prefix and has non-whitespace content after).
func IsAIQuery(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, aiPrefix) {
		return false
	}
	// Check for non-empty content after the prefix
	after := strings.TrimSpace(trimmed[len(aiPrefix):])
	return after != ""
}
