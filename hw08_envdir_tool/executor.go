package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	for _, arg := range cmd {
		if strings.ContainsAny(arg, "&|;<>\n") {
			return 1
		}
	}

	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	osEnv := os.Environ()
	newEnv := make([]string, 0, len(osEnv)+len(env))

	for _, envVar := range osEnv {
		//nolint:gocritic
		key := envVar[:strings.Index(envVar, "=")]
		if val, exists := env[key]; exists {
			if val.NeedRemove {
				continue
			}
			newEnv = append(newEnv, key+"="+val.Value)
		} else {
			newEnv = append(newEnv, envVar)
		}
	}

	for key, val := range env {
		if _, exists := os.LookupEnv(key); !exists {
			if !val.NeedRemove {
				newEnv = append(newEnv, key+"="+val.Value)
			}
		}
	}

	command.Env = newEnv

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
