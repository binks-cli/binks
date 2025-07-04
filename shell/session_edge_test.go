package shell

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/binks-cli/binks/internal/executor"
	"github.com/stretchr/testify/assert"
)

func TestSession_ChangeDir_EdgeCases(t *testing.T) {
	sess := NewSession()
	cwd := sess.Cwd()

	testCases := []struct {
		name      string
		input     string
		expectErr bool
		check     func(t *testing.T, sess *Session, err error)
	}{
		{"empty string (home)", "", false, func(t *testing.T, sess *Session, err error) {
			home, _ := os.UserHomeDir()
			assert.Equal(t, home, sess.Cwd())
		}},
		{"~ (home)", "~", false, func(t *testing.T, sess *Session, err error) {
			home, _ := os.UserHomeDir()
			assert.Equal(t, home, sess.Cwd())
		}},
		{"non-existent dir", "/no/such/dir/shouldexist", true, func(t *testing.T, sess *Session, err error) {
			assert.Error(t, err)
			assert.Equal(t, cwd, sess.Cwd(), "cwd should not change on error")
		}},
		{"relative valid dir", "..", false, func(t *testing.T, sess *Session, err error) {
			assert.NotEqual(t, cwd, sess.Cwd())
		}},
		{"path with ~ prefix", "~/", false, func(t *testing.T, sess *Session, err error) {
			home, _ := os.UserHomeDir()
			assert.True(t, strings.HasPrefix(sess.Cwd(), home))
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sess.ChangeDir(tc.input)
			tc.check(t, sess, err)
		})
	}
}

func TestSession_RunCommand_EdgeCases(t *testing.T) {
	sess := NewSession()

	// Replace executor with a mock for error simulation
	mock := &executor.MockExecutorTestify{}
	mock.On("RunCommand", "").Return("", nil)
	mock.On("RunCommand", "failing").Return("", errors.New("fail"))
	mock.On("RunCommand", strings.Repeat("a", 10000)).Return("ok", nil)
	sess.Executor = mock

	t.Run("empty command", func(t *testing.T) {
		output, err := sess.RunCommand("")
		assert.NoError(t, err)
		assert.Equal(t, "", output)
	})

	t.Run("failing command", func(t *testing.T) {
		_, err := sess.RunCommand("failing")
		assert.Error(t, err)
		assert.EqualError(t, err, "fail")
	})

	t.Run("very long command", func(t *testing.T) {
		output, err := sess.RunCommand(strings.Repeat("a", 10000))
		assert.NoError(t, err)
		assert.Equal(t, "ok", output)
	})

	mock.AssertExpectations(t)
}

func TestSession_ChangeDir_NoPanic(t *testing.T) {
	sess := NewSession()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ChangeDir panicked: %v", r)
		}
	}()
	_ = sess.ChangeDir("")
	_ = sess.ChangeDir("~")
	_ = sess.ChangeDir("/no/such/dir/shouldexist")
}
