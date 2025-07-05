package executor

// isInteractiveCommand returns true if the command is known to require an attached terminal (interactive mode)
func isInteractiveCommand(cmd string) bool {
	// List of common interactive commands
	interactive := []string{"vim", "nano", "less", "more", "man", "ssh", "top", "htop", "nvim", "vi"}
	for _, name := range interactive {
		if len(cmd) >= len(name) && cmd[:len(name)] == name {
			return true
		}
	}
	return false
}
