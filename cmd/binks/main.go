package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/binks-cli/binks/internal/executor"
	"github.com/binks-cli/binks/shell"
	"github.com/kballard/go-shellquote"
	"golang.org/x/term"
)

func enableAltScreen() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if _, err := fmt.Fprint(os.Stdout, "\x1b[?1049h"); err != nil {
			fmt.Fprintln(os.Stderr, "failed to enable alt screen:", err)
		}
		if err := os.Stdout.Sync(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to sync stdout:", err)
		}
	}
}

func disableAltScreen() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if _, err := fmt.Fprint(os.Stdout, "\x1b[?1049l"); err != nil {
			fmt.Fprintln(os.Stderr, "failed to disable alt screen:", err)
		}
		if err := os.Stdout.Sync(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to sync stdout:", err)
		}
	}
}

func main() {
	altScreen := os.Getenv("BINKS_ALT_SCREEN") == "1"
	if altScreen {
		enableAltScreen()
		// Ensure alt screen is disabled on SIGINT/SIGTERM
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			disableAltScreen()
			os.Exit(1)
		}()
		// Do not use defer here, as os.Exit will prevent it from running
	}
	if len(os.Args) < 2 {
		// Start interactive REPL mode
		sess := shell.NewSession()
		err := shell.RunREPL(sess)
		if err != nil {
			fmt.Fprint(os.Stderr, shell.ErrorMessage(err))
			os.Exit(1)
		}
		return
	}

	// Properly quote and join all arguments after the program name to form the command
	command := shellquote.Join(os.Args[1:]...)

	exec := executor.NewBashExecutor()
	output, err := exec.RunCommand(command)

	if err != nil {
		fmt.Fprint(os.Stderr, shell.ErrorMessage(err))
		os.Exit(1)
	}

	fmt.Print(output)
}
