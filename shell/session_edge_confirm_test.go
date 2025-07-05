package shell

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseAISuggestion(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		explanation string
		command     string
	}{
		{
			name:        "explanation and code block",
			input:       "Here is what you should do:\n```sh\ngit pull && make test\n```",
			explanation: "Here is what you should do:",
			command:     "git pull && make test",
		},
		{
			name:        "no code block",
			input:       "Just run the tests manually.",
			explanation: "Just run the tests manually.",
			command:     "",
		},
		{
			name:        "code block with bash prefix",
			input:       "Do this:\n```bash\nls -la\n```",
			explanation: "Do this:",
			command:     "ls -la",
		},
		{
			name:        "code block with no language",
			input:       "```\necho hi\n```",
			explanation: "",
			command:     "echo hi",
		},
		{
			name:        "multiline explanation and command",
			input:       "First, update:\nThen run:\n```sh\ngit pull && make test\n```",
			explanation: "First, update:\nThen run:",
			command:     "git pull && make test",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			exp, cmd := parseAISuggestion(tc.input)
			assert.Equal(t, tc.explanation, exp)
			assert.Equal(t, tc.command, cmd)
		})
	}
}

func TestSession_ExecuteLine_AIWithCodeBlock(t *testing.T) {
	s := &Session{
		Executor: nil,
		Agent: agentFuncMock(func(prompt string) (string, error) {
			return "Explanation here.\n```sh\necho dummy\n```", nil
		}),
		cwd: ".",
	}
	resp, err := s.ExecuteLine(">> do something")
	assert.NoError(t, err)
	assert.Equal(t, "[AI]", resp)
	assert.NotNil(t, s.pendingSuggestion)
	assert.Equal(t, "echo dummy", s.pendingSuggestion.command)
	assert.Equal(t, "Explanation here.", s.pendingSuggestion.explanation)
}

func TestSession_ExecuteLine_AINoCodeBlock(t *testing.T) {
	s := &Session{
		Executor: nil,
		Agent: agentFuncMock(func(prompt string) (string, error) {
			return "Just explain, no command.", nil
		}),
		cwd: ".",
	}
	resp, err := s.ExecuteLine(">> explain only")
	assert.NoError(t, err)
	assert.Equal(t, "[AI] Just explain, no command.", resp)
	assert.Nil(t, s.pendingSuggestion)
}

type agentFuncMock func(string) (string, error)

func (f agentFuncMock) Respond(prompt string) (string, error) {
	return f(prompt)
}
