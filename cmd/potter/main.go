package main

import (
	"os"

	"github.com/JulienBreux/potter/internal/potter/command"
)

func main() {
	cmd := command.New(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		command.PrintError(os.Stderr, err)
	}
}
