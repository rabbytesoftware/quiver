package logger

import (
	"fmt"
	"os"
)

func CompressFile(originalFile *os.File, compressFolderPath, level string) error{
	// Create compressed file with name format
	// "level-compressed-timestamp.txt.gz"
	createdFile, createErr := CreateLogFile(compressFolderPath, level, true)

	if createErr != nil {
		return fmt.Errorf("could not create compressed file: %w", createErr)
	}

	defer createdFile.Close()

	// Create gzip writer
	compressedFile := CompressToGzipFile(createdFile)
	defer compressedFile.Close()

	// Copy original file content
	// to the compressed one
	copyErr := CopyFileContentToCompressedGzipFile(originalFile, compressedFile)

	if copyErr != nil {
		return fmt.Errorf("could not copy original file content to the compressed one: %w", copyErr)
	}

	// Delete original file
	deleteErr := DeleteFile(originalFile)

	if deleteErr != nil {
		return fmt.Errorf("could not delete the original file: %w", deleteErr)
	}

	return nil
}