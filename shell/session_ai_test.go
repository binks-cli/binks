package shell

import (
	"bytes"
	"errors"
	"testing"

	"github.com/binks-cli/binks/internal/agent"
	"github.com/stretchr/testify/assert"
)

func TestSession_AI_ConfirmFlow(t *testing.T) {
	var out bytes.Buffer
	sess := &Session{
		Executor: &mockExecutor{},
		Agent: agent.NewMockAgent(map[string]agent.AgentResult{
			"clean": {Output: "Here is what you should do:\n```sh\nrm -rf build/\n```", Err: nil},
		}),
		cwd: ".",
		Out: &out,
	}
	resp, err := sess.ExecuteLine(">> clean")
	assert.NoError(t, err)
	assert.Equal(t, "[AI]", resp)
	assert.NotNil(t, sess.pendingSuggestion)
	assert.Equal(t, "rm -rf build/", sess.pendingSuggestion.command)
	assert.Contains(t, out.String(), "AI suggests: rm -rf build/")

	out.Reset()
	resp, err = sess.ExecuteLine("yes")
	assert.NoError(t, err)
	assert.Contains(t, resp, "executed: rm -rf build/")
	assert.Nil(t, sess.pendingSuggestion)
	assert.Contains(t, out.String(), "executed: rm -rf build/")

	// Decline scenario
	out.Reset()
	sess.pendingSuggestion = &PendingSuggestion{command: "echo hi"}
	resp, err = sess.ExecuteLine("no")
	assert.NoError(t, err)
	assert.Contains(t, resp, "Cancelled")
	assert.Nil(t, sess.pendingSuggestion)
	assert.Contains(t, out.String(), "Cancelled")
}

func TestSession_AI_ErrorAndNoCommand(t *testing.T) {
	var out, errBuf bytes.Buffer
	sess := &Session{
		Executor: &mockExecutor{},
		Agent: agent.NewMockAgent(map[string]agent.AgentResult{
			"fail":  {Output: "", Err: errors.New("API timeout")},
			"plain": {Output: "I cannot help with that", Err: nil},
		}),
		cwd: ".",
		Out: &out,
		Err: &errBuf,
	}
	resp, err := sess.ExecuteLine(">> fail")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API timeout")
	assert.Nil(t, sess.pendingSuggestion)
	assert.Contains(t, errBuf.String(), "[AI] error:")

	out.Reset()
	errBuf.Reset()
	resp, err = sess.ExecuteLine(">> plain")
	assert.NoError(t, err)
	assert.Contains(t, resp, "I cannot help with that")
	assert.Nil(t, sess.pendingSuggestion)
	assert.Contains(t, out.String(), "I cannot help with that")
}

func TestSession_AI_YesWithoutPending(t *testing.T) {
	var out bytes.Buffer
	exec := &mockExecutor{}
	sess := &Session{
		Executor: exec,
		Agent: agent.NewMockAgent(nil, agent.AgentResult{Output: "", Err: nil}),
		cwd: ".",
		Out: &out,
	}
	resp, err := sess.ExecuteLine("yes")
	assert.NoError(t, err)
	assert.Contains(t, resp, "executed: yes")
	assert.Equal(t, 1, exec.calls)
	assert.Contains(t, out.String(), "executed: yes")
}
