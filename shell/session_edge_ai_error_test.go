package shell

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test various error and edge cases for AI error handling and confirmation
func TestSession_ExecuteLine_AIErrorsAndEdgeCases(t *testing.T) {
	s := &Session{
		Executor: nil,
		Agent: agentFuncMock(func(prompt string) (string, error) {
			if prompt == "timeout" {
				return "", errors.New("AI request timed out")
			}
			if prompt == "apierror" {
				return "", errors.New("OpenAI API error: invalid_api_key")
			}
			if prompt == "network" {
				return "", errors.New("AI error: network unreachable")
			}
			if prompt == "malformed" {
				return "", errors.New("AI error: failed to parse response")
			}
			return "No code block here.", nil
		}),
		cwd: ".",
	}

	resp, err := s.ExecuteLine(">> timeout")
	assert.Error(t, err)
	assert.Contains(t, resp, "AI request timed out")
	assert.Nil(t, s.pendingSuggestion)

	resp, err = s.ExecuteLine(">> apierror")
	assert.Error(t, err)
	assert.Contains(t, resp, "OpenAI API error")
	assert.Nil(t, s.pendingSuggestion)

	resp, err = s.ExecuteLine(">> network")
	assert.Error(t, err)
	assert.Contains(t, resp, "AI error: network unreachable")
	assert.Nil(t, s.pendingSuggestion)

	resp, err = s.ExecuteLine(">> malformed")
	assert.Error(t, err)
	assert.Contains(t, resp, "failed to parse response")
	assert.Nil(t, s.pendingSuggestion)

	resp, err = s.ExecuteLine(">> no code block")
	assert.NoError(t, err)
	assert.Contains(t, resp, "No code block here.")
	assert.Nil(t, s.pendingSuggestion)
}

func TestSession_ExecuteLine_ConfirmationInputVariants(t *testing.T) {
	s := &Session{
		Executor: &mockExecutor{},
		Agent: agentFuncMock(func(prompt string) (string, error) {
			return "Explanation\n```sh\necho hi\n```", nil
		}),
		cwd: ".",
	}
	// Trigger pending suggestion
	resp, err := s.ExecuteLine(">> suggest")
	assert.NoError(t, err)
	assert.Equal(t, "[AI]", resp)
	assert.NotNil(t, s.pendingSuggestion)

	// Accept with 'y'
	s.pendingSuggestion = &PendingSuggestion{command: "echo hi"}
	resp, err = s.ExecuteLine("y")
	assert.NoError(t, err)
	assert.Nil(t, s.pendingSuggestion)

	// Accept with 'yes'
	s.pendingSuggestion = &PendingSuggestion{command: "echo hi"}
	_, err = s.ExecuteLine("yes")
	assert.NoError(t, err)
	assert.Nil(t, s.pendingSuggestion)

	// Decline with 'n', 'no', 'abc', '', 'Yess', 'YES!'
	for _, input := range []string{"n", "no", "abc", "", "Yess", "YES!"} {
		s.pendingSuggestion = &PendingSuggestion{command: "echo hi"}
		resp, err = s.ExecuteLine(input)
		assert.NoError(t, err)
		assert.Equal(t, "[AI] Cancelled.", resp)
		assert.Nil(t, s.pendingSuggestion)
	}
}
