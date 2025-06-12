package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if err := validateArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := copyFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func validateArgs() error {
	if from == "" || to == "" {
		return fmt.Errorf("both 'from' and 'to' paths must be specified")
	}

	if offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}

	return nil
}

func copyFile() error {
	srcFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("source file is empty")
	}

	if offset > fileInfo.Size() {
		return fmt.Errorf("offset exceeds file size")
	}

	if _, err := srcFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to offset: %w", err)
	}

	destFile, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	var bytesToCopy int64
	if limit == 0 {
		bytesToCopy = fileInfo.Size() - offset
	} else {
		bytesToCopy = min(limit, fileInfo.Size()-offset)
	}

	buf := make([]byte, 64*1024) // 32KB buffer
	var totalCopied int64
	lastProgress := -1

	for totalCopied < bytesToCopy {
		remaining := bytesToCopy - totalCopied
		if int64(len(buf)) > remaining {
			buf = buf[:remaining]
		}

		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading from source: %w", err)
		}
		if n == 0 {
			break
		}

		if _, err := destFile.Write(buf[:n]); err != nil {
			return fmt.Errorf("error writing to destination: %w", err)
		}

		totalCopied += int64(n)
		progress := int(float64(totalCopied) / float64(bytesToCopy) * 100)

		if progress != lastProgress {
			fmt.Printf("\rProgress: %d%%", progress)
			lastProgress = progress
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("\nCopy completed successfully!")
	return nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
