package shell

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/binks-cli/binks/internal/executor"
	"github.com/stretchr/testify/assert"
)

func TestSession_ChangeDir_EdgeCases(t *testing.T) {
	tmp := t.TempDir()
	_ = os.Chdir(tmp)

	testCases := []struct {
		name      string
		input     string
		expectErr bool
		check     func(_ *testing.T, sess *Session, startCwd string, _ error)
	}{
		{"empty string (home)", "", false, func(_ *testing.T, sess *Session, startCwd string, _ error) {
			home, _ := os.UserHomeDir()
			assert.Equal(t, home, sess.Cwd())
		}},
		{"~ (home)", "~", false, func(_ *testing.T, sess *Session, startCwd string, _ error) {
			home, _ := os.UserHomeDir()
			assert.Equal(t, home, sess.Cwd())
		}},
		{"non-existent dir", "/no/such/dir/shouldexist", true, func(_ *testing.T, sess *Session, startCwd string, err error) {
			assert.Error(t, err)
			assert.Equal(t, startCwd, sess.Cwd(), "cwd should not change on error")
		}},
		{"relative valid dir", "..", false, func(_ *testing.T, sess *Session, startCwd string, _ error) {
			assert.NotEqual(t, startCwd, sess.Cwd())
		}},
		{"path with ~ prefix", "~/", false, func(_ *testing.T, sess *Session, startCwd string, _ error) {
			home, _ := os.UserHomeDir()
			assert.True(t, strings.HasPrefix(sess.Cwd(), home))
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sess := NewSession()
			startCwd := sess.Cwd()
			err := sess.ChangeDir(tc.input)
			tc.check(t, sess, startCwd, err)
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

func TestRunCommand_DoesNotAffectSessionCwd(t *testing.T) {
	tmp := t.TempDir()
	sess := NewSession()
	startCwd := sess.Cwd()

	// Run an external cd command (should not affect session cwd)
	_, err := sess.RunCommand("cd " + tmp)
	assert.NoError(t, err)
	cwdEval, _ := filepath.EvalSymlinks(sess.Cwd())
	startEval, _ := filepath.EvalSymlinks(startCwd)
	assert.Equal(t, startEval, cwdEval, "Session cwd should not change after external cd command")

	// Now use built-in cd (should change session cwd)
	err = sess.ChangeDir(tmp)
	assert.NoError(t, err)
	cwdEval, _ = filepath.EvalSymlinks(sess.Cwd())
	tmpEval, _ := filepath.EvalSymlinks(tmp)
	assert.Equal(t, tmpEval, cwdEval, "Session cwd should change after built-in cd")
}

func TestSession_ChangeDir_PlatformSpecific(t *testing.T) {
	// This test checks platform-specific path handling for cd
	// On Windows, test drive letter; on Unix, test root
	if os.PathSeparator == '\\' {
		// Windows
		start := NewSession().Cwd()
		// Try to cd to C:\ (skip if not present)
		if _, err := os.Stat("C:\\"); err == nil {
			sess := NewSession()
			err := sess.ChangeDir("C:\\")
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(sess.Cwd(), "C:"))
			// cd .. from root should stay at root
			err = sess.ChangeDir("..")
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(sess.Cwd(), "C:"))
		}
		_ = start // avoid unused
	} else {
		// Unix
		sess := NewSession()
		err := sess.ChangeDir("/")
		assert.NoError(t, err)
		assert.Equal(t, "/", sess.Cwd())
		// cd .. from root should stay at root
		err = sess.ChangeDir("..")
		assert.NoError(t, err)
		assert.Equal(t, "/", sess.Cwd())
	}
}
