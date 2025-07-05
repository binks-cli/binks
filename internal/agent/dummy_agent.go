package agent

import "fmt"

// DummyAgent is a stub implementation of the Agent interface.
type DummyAgent struct{}

// Respond returns a fixed response for development/testing.
func (d *DummyAgent) Respond(prompt string) (string, error) {
	return fmt.Sprintf("Echo: %s", prompt), nil
}
