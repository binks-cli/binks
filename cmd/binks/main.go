package main

import (
	"fmt"
	"os"

	"github.com/binks-cli/binks/internal/executor"
	"github.com/binks-cli/binks/shell"
	"github.com/kballard/go-shellquote"
)

func main() {
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
