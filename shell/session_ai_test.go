package shell

import (
	"errors"
	"testing"

	"github.com/binks-cli/binks/internal/agent"
	"github.com/stretchr/testify/assert"
)

func TestSession_AI_ConfirmFlow(t *testing.T) {
	sess := &Session{
		Executor: &mockExecutor{},
		Agent: agent.NewMockAgent(map[string]agent.AgentResult{
			"clean": {Output: "Here is what you should do:\n```sh\nrm -rf build/\n```", Err: nil},
		}),
		cwd: ".",
	}
	resp, err := sess.ExecuteLine(">> clean")
	assert.NoError(t, err)
	assert.Equal(t, "[AI]", resp)
	assert.NotNil(t, sess.pendingSuggestion)
	assert.Equal(t, "rm -rf build/", sess.pendingSuggestion.command)

	resp, err = sess.ExecuteLine("yes")
	assert.NoError(t, err)
	assert.Contains(t, resp, "executed: rm -rf build/")
	assert.Nil(t, sess.pendingSuggestion)

	// Decline scenario
	sess.pendingSuggestion = &PendingSuggestion{command: "echo hi"}
	resp, err = sess.ExecuteLine("no")
	assert.NoError(t, err)
	assert.Contains(t, resp, "Cancelled")
	assert.Nil(t, sess.pendingSuggestion)
}

func TestSession_AI_ErrorAndNoCommand(t *testing.T) {
	sess := &Session{
		Executor: &mockExecutor{},
		Agent: agent.NewMockAgent(map[string]agent.AgentResult{
			"fail":  {Output: "", Err: errors.New("API timeout")},
			"plain": {Output: "I cannot help with that", Err: nil},
		}),
		cwd: ".",
	}
	resp, err := sess.ExecuteLine(">> fail")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API timeout")
	assert.Nil(t, sess.pendingSuggestion)

	resp, err = sess.ExecuteLine(">> plain")
	assert.NoError(t, err)
	assert.Contains(t, resp, "I cannot help with that")
	assert.Nil(t, sess.pendingSuggestion)
}

func TestSession_AI_YesWithoutPending(t *testing.T) {
	exec := &mockExecutor{}
	sess := &Session{
		Executor: exec,
		Agent:    agent.NewMockAgent(nil, agent.AgentResult{Output: "", Err: nil}),
		cwd:      ".",
	}
	resp, err := sess.ExecuteLine("yes")
	assert.NoError(t, err)
	assert.Contains(t, resp, "executed: yes")
	assert.Equal(t, 1, exec.calls)
}
