package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		srcContent  string
		offset      int64
		limit       int64
		wantContent string
		wantError   bool
	}{
		{
			name:        "full copy",
			srcContent:  "test data",
			offset:      0,
			limit:       0,
			wantContent: "test data",
			wantError:   false,
		},
		{
			name:        "copy with offset",
			srcContent:  "test data",
			offset:      5,
			limit:       0,
			wantContent: "data",
			wantError:   false,
		},
		{
			name:        "copy with limit",
			srcContent:  "test data",
			offset:      0,
			limit:       4,
			wantContent: "test",
			wantError:   false,
		},
		{
			name:        "offset exceeds file size",
			srcContent:  "test",
			offset:      10,
			limit:       0,
			wantContent: "",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcPath := filepath.Join(tmpDir, "src.txt")
			err := os.WriteFile(srcPath, []byte(tt.srcContent), 0644)
			if err != nil {
				t.Fatal(err)
			}

			dstPath := filepath.Join(tmpDir, "dst.txt")

			from = srcPath
			to = dstPath
			offset = tt.offset
			limit = tt.limit

			err = copyFile()

			if (err != nil) != tt.wantError {
				t.Errorf("copyFile() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				got, err := os.ReadFile(dstPath)
				if err != nil {
					t.Fatal(err)
				}

				if string(got) != tt.wantContent {
					t.Errorf("copyFile() = %v, want %v", string(got), tt.wantContent)
				}
			}
		})
	}
}

func TestMainArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "src.txt")
	dstPath := filepath.Join(tmpDir, "dst.txt")
	err := os.WriteFile(srcPath, []byte("test data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Helper function to capture stdout
	runWithCapture := func(args []string) string {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = args
		main()

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}

	t.Run("success case", func(t *testing.T) {
		output := runWithCapture([]string{"cmd", "-from", srcPath, "-to", dstPath})
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			t.Error("output file was not created")
		}
		if output == "" {
			t.Error("expected progress output")
		}
	})

	t.Run("error case", func(t *testing.T) {
		exitCalled := false
		oldExit := osExit
		defer func() { osExit = oldExit }()
		osExit = func(code int) { exitCalled = true }

		runWithCapture([]string{"cmd", "-from", "nonexistent.txt", "-to", dstPath})

		if !exitCalled {
			t.Error("expected program to exit with error")
		}
	})
}

// Mock for os.Exit
var osExit = func(code int) {
	panic(code)
}
