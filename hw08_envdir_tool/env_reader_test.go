package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	// Создаем временную директорию для тестов
	dir, err := os.MkdirTemp("", "envdir_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Создаем тестовые файлы
	tests := []struct {
		filename string
		content  string
		expected EnvValue
	}{
		{"FOO", "123", EnvValue{Value: "123"}},
		{"BAR", "value \t ", EnvValue{Value: "value"}},
		{"EMPTY", "", EnvValue{NeedRemove: true}},
		{"NULLS", "hello\x00world", EnvValue{Value: "hello\nworld"}},
		{"INVALID=NAME", "test", EnvValue{}}, // файл с = в имени игнорируется
	}

	for _, tt := range tests {
		if tt.filename == "INVALID=NAME" {
			continue
		}
		err := os.WriteFile(filepath.Join(dir, tt.filename), []byte(tt.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", tt.filename, err)
		}
	}

	// Тестируем ReadDir
	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	// Проверяем результаты
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if tt.filename == "INVALID=NAME" {
				if _, exists := env[tt.filename]; exists {
					t.Errorf("File with invalid name %s was not ignored", tt.filename)
				}
				return
			}
			got, exists := env[tt.filename]
			if !exists {
				t.Errorf("Variable %s not found in env", tt.filename)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("For %s, got %+v, want %+v", tt.filename, got, tt.expected)
			}
		})
	}

	// Проверяем обработку несуществующей директории
	_, err = ReadDir("/nonexistent/dir")
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
}
