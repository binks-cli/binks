package shell

import (
	"strings"

	"github.com/binks-cli/binks/internal/agent"
)

// ExecuteLine dispatches input to the shell executor or the Agent, depending on isAIQuery.
func (s *Session) ExecuteLine(line string) (string, error) {
	if agent.IsAIQuery(line) && s.Agent != nil {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ">>") {
			trimmed = strings.TrimSpace(trimmed[2:])
		}
		resp, err := s.Agent.Respond(trimmed)
		if err != nil {
			return "[AI] error: " + err.Error(), err
		}
		// Parse AI response for code block (shell command)
		explanation, command := parseAISuggestion(resp)
		if command != "" {
			s.pendingSuggestion = &PendingSuggestion{
				explanation: explanation,
				command:     command,
				raw:         resp,
				confirmed:   false,
				declined:    false,
			}
			return "[AI]", nil // Signal to REPL to prompt for confirmation
		}
		return "[AI] " + resp, nil
	}
	return s.RunCommand(line)
}

// parseAISuggestion extracts explanation and the first shell command code block from AI response.
func parseAISuggestion(resp string) (explanation, command string) {
	resp = strings.ReplaceAll(resp, "\r\n", "\n")
	parts := strings.Split(resp, "```")
	if len(parts) < 3 {
		return resp, "" // No code block found
	}
	// parts[0]: explanation, parts[1]: language (optional), parts[2]: code
	explanation = strings.TrimSpace(parts[0])
	code := strings.TrimSpace(parts[2])
	// If the code block starts with 'sh' or 'bash', skip that line
	if strings.HasPrefix(code, "sh\n") {
		code = strings.TrimPrefix(code, "sh\n")
	} else if strings.HasPrefix(code, "bash\n") {
		code = strings.TrimPrefix(code, "bash\n")
	}
	return explanation, strings.TrimSpace(code)
}
