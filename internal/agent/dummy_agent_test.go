package agent

import (
	"testing"
)

func TestDummyAgent_Respond(t *testing.T) {
	agent := &DummyAgent{}
	prompt := "Hello, Agent!"
	response, err := agent.Respond(prompt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	want := "Echo: Hello, Agent!"
	if response != want {
		t.Errorf("Expected response %q, got %q", want, response)
	}
}
