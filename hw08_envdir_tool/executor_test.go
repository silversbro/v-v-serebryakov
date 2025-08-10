package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCmd(t *testing.T) {
	// Создаем временную директорию для тестов
	dir, err := os.MkdirTemp("", "envdir_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Создаем тестовый скрипт echo.sh
	echoScript := `#!/bin/bash
echo -e "FOO=($FOO)\nBAR=($BAR)\nADDED=($ADDED)"
`
	scriptPath := filepath.Join(dir, "echo.sh")
	err = os.WriteFile(scriptPath, []byte(echoScript), 0600)
	if err != nil {
		t.Fatalf("Failed to create echo.sh: %v", err)
	}

	// Создаем файлы окружения
	os.WriteFile(filepath.Join(dir, "FOO"), []byte("foovalue"), 0644)
	os.WriteFile(filepath.Join(dir, "BAR"), []byte("barvalue"), 0644)
	os.WriteFile(filepath.Join(dir, "EMPTY"), []byte(""), 0644)

	// Устанавливаем начальные переменные окружения
	os.Setenv("ADDED", "original")
	defer os.Unsetenv("ADDED")

	// Читаем окружение
	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	// Тестируем RunCmd
	t.Run("Run command with modified env", func(t *testing.T) {
		cmd := []string{scriptPath}
		returnCode := RunCmd(cmd, env)
		if returnCode != 0 {
			t.Errorf("Expected return code 0, got %d", returnCode)
		}
	})

	// Проверяем удаление переменной
	t.Run("Remove env variable", func(t *testing.T) {
		env := Environment{"EMPTY": EnvValue{NeedRemove: true}}
		os.Setenv("EMPTY", "shouldberemoved")
		defer os.Unsetenv("EMPTY")

		cmd := []string{scriptPath}
		returnCode := RunCmd(cmd, env)
		if returnCode != 0 {
			t.Errorf("Expected return code 0, got %d", returnCode)
		}
	})

	// Проверяем пустую команду
	t.Run("Empty command", func(t *testing.T) {
		cmd := []string{}
		returnCode := RunCmd(cmd, env)
		if returnCode != 1 {
			t.Errorf("Expected return code 1 for empty command, got %d", returnCode)
		}
	})

	// Проверяем несуществующую команду
	t.Run("Nonexistent command", func(t *testing.T) {
		cmd := []string{"/nonexistent/command"}
		returnCode := RunCmd(cmd, env)
		if returnCode == 0 {
			t.Error("Expected non-zero return code for nonexistent command")
		}
	})
}
