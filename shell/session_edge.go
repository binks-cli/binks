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
		return "[AI] " + resp, nil
	}
	return s.RunCommand(line)
}
