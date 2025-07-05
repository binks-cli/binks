package shell

// mockExecutor is a test double for Executor, used in multiple test files.
type mockExecutor struct {
	lastCmd string
	calls   int
	fail    bool
	resp    string
	err     error
}

func (m *mockExecutor) RunCommand(cmd string) (string, error) {
	m.lastCmd = cmd
	m.calls++
	if m.fail {
		return "", m.err
	}
	if m.resp != "" {
		return m.resp, m.err
	}
	return "executed: " + cmd, nil
}

func (m *mockExecutor) RunCommandWithDir(cmd, dir string) (string, error) {
	return m.RunCommand(cmd)
}
