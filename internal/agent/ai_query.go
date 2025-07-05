package agent

import "strings"

// AIPrefix is the user-facing prefix for AI queries.
const AIPrefix = ">>"

// IsAIQuery returns true if the line is an AI query (starts with the AI prefix and has non-whitespace content after).
func IsAIQuery(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, AIPrefix) {
		return false
	}
	// Check for non-empty content after the prefix
	after := strings.TrimSpace(trimmed[len(AIPrefix):])
	return after != ""
}
