package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
)

var (
	// ErrFileGetInfo возвращается, не получилось получить информацию.
	ErrFileGetInfo = errors.New("error file get info")
	// ErrCopyFileToSelf возвращается, не получилось получить информацию.
	ErrCopyFileToSelf = errors.New("error copy file to self")
	// ErrOpeningFile возвращается, когда читаемый файл не поддерживается.
	ErrOpeningFile = errors.New("error opening file")
	// ErrReadingFromFile возвращается, когда ошибка при чтении файла.
	ErrReadingFromFile = errors.New("error reading from file")
	// ErrOffsetExceedsFileSize возвращается, когда отступ больше чем размер файла.
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	// ErrCreateDestinationFile возвращается, когда не возможно создать файл для записи.
	ErrCreateDestinationFile = errors.New("failed to create destination file")
	// ErrWriteDestinationFile возвращается, когда не возможно записать в файл.
	ErrWriteDestinationFile = errors.New("error writing to destination")
)

type FileInfo = fs.FileInfo

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return getError(ErrCopyFileToSelf, nil)
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsPermission(err) {
			return getError(ErrReadingFromFile, err)
		}
		return getError(ErrOpeningFile, err)
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

	// Проверяем, является ли файл устройством или специальным файлом
	isSpecial := (fileInfo.Mode()&os.ModeDevice != 0) ||
		(fileInfo.Mode()&os.ModeNamedPipe != 0) ||
		(fileInfo.Mode()&os.ModeSocket != 0) ||
		(fileInfo.Mode()&os.ModeCharDevice != 0)

	if isSpecial {
		return copySpecialFile(srcFile, toPath, offset, limit)
	}
	return copyRegularFile(srcFile, fileInfo, toPath, offset, limit)
}

func copyRegularFile(srcFile *os.File, fileInfo os.FileInfo, toPath string, offset, limit int64) error {
	if offset > fileInfo.Size() {
		return getError(ErrOffsetExceedsFileSize, nil)
	}

	if _, err := srcFile.Seek(offset, io.SeekStart); err != nil {
		return getError(ErrReadingFromFile, err)
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return getError(ErrCreateDestinationFile, err)
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			log.Print("Error closing file: ", err)
		}
	}()

	bytesToCopy := fileInfo.Size() - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	bufSize := getBufferSize(bytesToCopy)
	buf := make([]byte, bufSize)
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
	}

	fmt.Println("\nCopy completed successfully!")
	return nil
}

func copySpecialFile(srcFile *os.File, toPath string, offset, limit int64) error {
	if offset > 0 {
		if _, err := srcFile.Seek(offset, io.SeekStart); err != nil {
			return getError(ErrReadingFromFile, err)
		}
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return getError(ErrCreateDestinationFile, err)
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			log.Print("Error closing file: ", err)
		}
	}()

	const bufSize = 32 * 1024 // Фиксированный размер буфера
	buf := make([]byte, bufSize)
	var totalCopied int64

	for {
		if limit > 0 && totalCopied >= limit {
			break
		}

		readSize := bufSize
		if limit > 0 {
			remaining := limit - totalCopied
			if remaining < int64(bufSize) {
				readSize = int(remaining)
			}
		}

		n, err := io.ReadFull(srcFile, buf[:readSize])
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return getError(ErrReadingFromFile, err)
		}

		if n == 0 {
			break
		}

		if _, err := destFile.Write(buf[:n]); err != nil {
			return getError(ErrWriteDestinationFile, err)
		}

		totalCopied += int64(n)
		if limit > 0 {
			progress := int(float64(totalCopied) / float64(limit) * 100)
			fmt.Printf("\rProgress: %d%%", progress)
		} else {
			fmt.Printf("\rCopied: %d bytes", totalCopied)
		}
	}

	fmt.Println("\nCopy completed successfully!")
	return nil
}

func getBufferSize(fileSize int64) int {
	switch {
	case fileSize > 100*1024*1024:
		return 10 * 1024 * 1024
	case fileSize > 1024*1024:
		return 512 * 1024
	case fileSize < 1024*1024:
		return int(fileSize)
	default:
		return 128 * 1024
	}
}

func getError(errors error, errSys error) error {
	if errSys == nil {
		return fmt.Errorf("%w", errors)
	}

	return fmt.Errorf("%w: %v", errors, errSys) //nolint:errorlint
}
