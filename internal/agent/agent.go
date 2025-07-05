package agent

// Agent is an interface for responding to prompts.
type Agent interface {
	Respond(prompt string) (string, error)
}

// AgentFunc allows using a function as an Agent for testing.
type AgentFunc func(string) (string, error)

func (f AgentFunc) Respond(prompt string) (string, error) {
	return f(prompt)
}
