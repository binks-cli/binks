package agent

import "errors"

// AgentResult represents a mock response from the agent.
type AgentResult struct {
	Output string
	Err    error
}

// MockAgent implements Agent and returns predefined responses for prompts.
type MockAgent struct {
	Responses map[string]AgentResult
	Default   AgentResult // fallback if prompt not found
}

func (m *MockAgent) Respond(prompt string) (string, error) {
	if m.Responses == nil {
		return m.Default.Output, m.Default.Err
	}
	if res, ok := m.Responses[prompt]; ok {
		return res.Output, res.Err
	}
	return m.Default.Output, m.Default.Err
}

// NewMockAgent creates a MockAgent with the given responses and optional default.
func NewMockAgent(responses map[string]AgentResult, defaultResult ...AgentResult) *MockAgent {
	var def AgentResult
	if len(defaultResult) > 0 {
		def = defaultResult[0]
	} else {
		def = AgentResult{"", errors.New("no mock response")}
	}
	return &MockAgent{Responses: responses, Default: def}
}
