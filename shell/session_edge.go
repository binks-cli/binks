package shell

import (
	"regexp"
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
	// Regex to match code block: ```[lang]?\n...\n```
	re := regexp.MustCompile("(?s)```(?:[a-zA-Z]+)?\\n(.*?)```")
	match := re.FindStringSubmatch(resp)
	if match != nil {
		cmd := strings.TrimSpace(match[1])
		// Remove the code block from the response for explanation
		explanation = strings.TrimSpace(re.ReplaceAllString(resp, ""))
		return explanation, cmd
	}
	return resp, ""
}
