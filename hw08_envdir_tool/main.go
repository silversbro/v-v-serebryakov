package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s <env_dir> <command> [args...]\n", os.Args[0])
		if err != nil {
			fmt.Printf("Fprintf: %v\n", err)
		}
		os.Exit(1)
	}

	envDir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error reading env directory: %v\n", err)
		if err != nil {
			fmt.Printf("Fprintf: %v\n", err)
		}
		os.Exit(1)
	}

	os.Exit(RunCmd(cmd, env))
}
