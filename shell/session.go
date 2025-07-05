package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/binks-cli/binks/internal/agent"
	"github.com/binks-cli/binks/internal/executor"
)

// Session represents the state of a shell session
type Session struct {
	Executor  executor.Executor
	Agent     agent.Agent // AI agent for handling AI queries
	cwd       string      // Current working directory
	AIEnabled bool        // Global AI mode toggle
	// Future fields for working directory, history, etc.
}

// NewSession creates a new shell session
func NewSession() *Session {
	wd, err := os.Getwd()
	if err != nil {
		wd = "." // fallback
	}
	var ag agent.Agent
	if os.Getenv("OPENAI_API_KEY") != "" {
		ag = agent.NewOpenAIAgent()
	} else {
		ag = &agent.DummyAgent{}
	}
	return &Session{
		Executor:  executor.NewBashExecutor(),
		Agent:     ag,
		cwd:       wd,
		AIEnabled: false, // Default to off
	}
}

// Cwd returns the current working directory for the session
func (s *Session) Cwd() string {
	return s.cwd
}

// ChangeDir changes the session's current working directory
func (s *Session) ChangeDir(path string) error {
	var target string
	switch {
	case path == "":
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = home
	case path == "~":
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = home
	case strings.HasPrefix(path, "~"):
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = filepath.Join(home, path[1:])
	default:
		target = path
	}
	if err := os.Chdir(target); err != nil {
		return err
	}
	abs, err := os.Getwd()
	if err != nil {
		return err
	}
	s.cwd = abs
	return nil
}

// RunCommand runs a command in the session's current working directory
func (s *Session) RunCommand(cmd string) (string, error) {
	if be, ok := s.Executor.(*executor.BashExecutor); ok {
		return be.RunCommandWithDir(cmd, s.cwd)
	}
	return s.Executor.RunCommand(cmd)
}
