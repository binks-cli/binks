package main

import (
	"fmt"
	"os"

	"github.com/binks-cli/binks/internal/executor"
	"github.com/kballard/go-shellquote"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: binks [command]")
		fmt.Println("Interactive mode coming in Stage 2")
		os.Exit(1)
	}

	// Properly quote and join all arguments after the program name to form the command
	command := shellquote.Join(os.Args[1:]...)
	
	exec := executor.NewBashExecutor()
	output, err := exec.RunCommand(command)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Print(output)
}
