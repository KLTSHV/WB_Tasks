package main

import (
	"cmd/internal/shell"
	"os"
)

func main() {
	sh := shell.New(shell.Options{
		Prompt: true,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	sh.Run()
}
