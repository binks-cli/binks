package agent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockAgent_Respond(t *testing.T) {
	responses := map[string]AgentResult{
		"foo": {Output: "echo hi", Err: nil},
		"bar": {Output: "I cannot help with that", Err: nil},
		"baz": {Output: "", Err: errors.New("API timeout")},
	}
	mock := NewMockAgent(responses, AgentResult{"default", nil})

	t.Run("returns mapped response", func(t *testing.T) {
		out, err := mock.Respond("foo")
		assert.NoError(t, err)
		assert.Equal(t, "echo hi", out)
	})
	t.Run("returns mapped error", func(t *testing.T) {
		out, err := mock.Respond("baz")
		assert.Error(t, err)
		assert.Equal(t, "API timeout", err.Error())
		assert.Equal(t, "", out)
	})
	t.Run("returns default for unknown", func(t *testing.T) {
		out, err := mock.Respond("unknown")
		assert.NoError(t, err)
		assert.Equal(t, "default", out)
	})
}
