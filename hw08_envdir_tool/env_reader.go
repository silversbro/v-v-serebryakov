package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if strings.Contains(filename, "=") {
			continue
		}

		filePath := filepath.Join(dir, filename)
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				if err == nil {
					err = fmt.Errorf("close error: %w", closeErr)
				}
			}
		}()

		fileInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if fileInfo.Size() == 0 {
			env[filename] = EnvValue{NeedRemove: true}
			continue
		}

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimRight(line, " \t")
			line = strings.ReplaceAll(line, "\x00", "\n")
			env[filename] = EnvValue{Value: line}
		} else {
			env[filename] = EnvValue{NeedRemove: true}
		}
	}

	return env, nil
}
