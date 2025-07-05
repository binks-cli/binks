package agent

// Agent is an interface for responding to prompts.
type Agent interface {
	Respond(prompt string) (string, error)
}
