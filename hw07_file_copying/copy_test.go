package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	// Создаем временную директорию для тестов
	testDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	tests := []struct {
		name        string
		srcContent  string
		offset      int64
		limit       int64
		expectError error
	}{
		{
			name:       "full copy small file",
			srcContent: "hello world",
			offset:     0,
			limit:      0,
		},
		{
			name:       "copy with offset",
			srcContent: "hello world",
			offset:     6,
			limit:      0,
		},
		{
			name:       "copy with limit",
			srcContent: "hello world",
			offset:     0,
			limit:      5,
		},
		{
			name:       "copy with offset and limit",
			srcContent: "hello world",
			offset:     2,
			limit:      5,
		},
		{
			name:        "empty file",
			srcContent:  "",
			offset:      0,
			limit:       0,
			expectError: ErrEmptyFile,
		},
		{
			name:        "offset exceeds file size",
			srcContent:  "hello",
			offset:      10,
			limit:       0,
			expectError: ErrOffsetExceedsFileSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем исходный файл
			srcFile, err := ioutil.TempFile(testDir, "src_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(srcFile.Name())

			if _, err := srcFile.WriteString(tt.srcContent); err != nil {
				t.Fatalf("Failed to write to src file: %v", err)
			}
			srcFile.Close()

			// Создаем файл назначения
			dstFile, err := ioutil.TempFile(testDir, "dst_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			dstFile.Close()
			defer os.Remove(dstFile.Name())

			// Выполняем копирование
			err = Copy(srcFile.Name(), dstFile.Name(), tt.offset, tt.limit)

			// Проверяем ошибки
			if tt.expectError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectError)
				} else if !errors.Is(err, tt.expectError) {
					t.Errorf("Expected error %v, got %v", tt.expectError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Проверяем содержимое скопированного файла
			expectedContent := calculateExpectedContent(tt.srcContent, tt.offset, tt.limit)
			actualContent, err := ioutil.ReadFile(dstFile.Name())
			if err != nil {
				t.Errorf("Failed to read dst file: %v", err)
				return
			}

			if string(actualContent) != expectedContent {
				t.Errorf("Expected content '%s', got '%s'", expectedContent, actualContent)
			}
		})
	}
}

func calculateExpectedContent(content string, offset, limit int64) string {
	if offset > int64(len(content)) {
		return ""
	}

	if limit == 0 {
		return content[offset:]
	}

	end := offset + limit
	if end > int64(len(content)) {
		end = int64(len(content))
	}

	return content[offset:end]
}

func TestCopyToNonexistentDir(t *testing.T) {
	testDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Создаем исходный файл
	srcFile, err := ioutil.TempFile(testDir, "src_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(srcFile.Name())

	if _, err := srcFile.WriteString("test content"); err != nil {
		t.Fatalf("Failed to write to src file: %v", err)
	}
	srcFile.Close()

	// Пытаемся скопировать в несуществующую директорию
	nonexistentPath := filepath.Join(testDir, "nonexistent_dir", "file.txt")
	err = Copy(srcFile.Name(), nonexistentPath, 0, 0)

	if err == nil {
		t.Error("Expected error when copying to nonexistent directory, got nil")
	}
}

func TestCopyFromNonexistentFile(t *testing.T) {
	testDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Создаем файл назначения
	dstFile, err := ioutil.TempFile(testDir, "dst_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dstFile.Close()
	defer os.Remove(dstFile.Name())

	// Пытаемся скопировать из несуществующего файла
	nonexistentPath := filepath.Join(testDir, "nonexistent_file.txt")
	err = Copy(nonexistentPath, dstFile.Name(), 0, 0)

	if err == nil {
		t.Error("Expected error when copying from nonexistent file, got nil")
	} else if !errors.Is(err, ErrUnsupportedFile) {
		t.Errorf("Expected error %v, got %v", ErrUnsupportedFile, err)
	}
}
