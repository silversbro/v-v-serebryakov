package hw07filecopying

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	// ErrEmptyFile возвращается, когда исходный файл пуст.
	ErrEmptyFile = errors.New("source file is empty")
	// ErrFileGetInfo возвращается, не получилось получить информацию.
	ErrFileGetInfo = errors.New("error file get info")
	// ErrUnsupportedFile возвращается, когда читаемый файл не поддерживается.
	ErrUnsupportedFile = errors.New("unsupported file")
	// ErrReadingFromFile возвращается, когда ошибка при чтении файла.
	ErrReadingFromFile = errors.New("error reading from file")
	// ErrOffsetExceedsFileSize возвращается, когда отступ больше чем размер файла.
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	// ErrCreateDestinationFile возвращается, когда не возможно создать файл для записи.
	ErrCreateDestinationFile = errors.New("failed to create destination file")
	// ErrWriteDestinationFile возвращается, когда не возможно записать в файл.
	ErrWriteDestinationFile = errors.New("error writing to destination")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath) // #nosec G304
	if err != nil {
		if os.IsPermission(err) {
			return getError(ErrReadingFromFile, err)
		}
		return getError(ErrUnsupportedFile, err)
	}

	defer func() {
		if err := srcFile.Close(); err != nil {
			log.Print("Error closing file: ", err)
		}
	}()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return getError(ErrFileGetInfo, err)
	}

	if fileInfo.Size() == 0 {
		return getError(ErrEmptyFile, nil)
	}

	if offset > fileInfo.Size() {
		return getError(ErrOffsetExceedsFileSize, nil)
	}

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return getError(ErrReadingFromFile, err)
	}

	destFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600) // #nosec G304
	if err != nil {
		return getError(ErrCreateDestinationFile, err)
	}

	defer func() {
		if err := destFile.Close(); err != nil {
			log.Print("Error closing file destination: ", err)
		}
	}()

	var bytesToCopy int64
	if limit == 0 {
		bytesToCopy = fileInfo.Size() - offset
	} else {
		bytesToCopy = minFinder(limit, fileInfo.Size()-offset)
	}

	buf := make([]byte, 64*1024)
	var totalCopied int64
	lastProgress := -1

	for totalCopied < bytesToCopy {
		remaining := bytesToCopy - totalCopied
		if int64(len(buf)) > remaining {
			buf = buf[:remaining]
		}

		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return getError(ErrReadingFromFile, err)
		}
		if n == 0 {
			break
		}

		if _, err := destFile.Write(buf[:n]); err != nil {
			return getError(ErrWriteDestinationFile, err)
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

func minFinder(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func getError(errors error, errSys error) error {
	if errSys == nil {
		return fmt.Errorf("%w", errors)
	}

	return fmt.Errorf("%s: %w", errors.Error(), errSys)
}
