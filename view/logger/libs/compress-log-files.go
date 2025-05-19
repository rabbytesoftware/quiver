package logger

import (
	"fmt"
	"os"
)

func CompressFile(originalFile *os.File, compressFolderPath, level string) error{
	// Create compressed file
	// "level-compressed-timestamp.txt.gz"
	createdFile, createErr := CreateFile(compressFolderPath, level, true)

	if createErr != nil {
		return fmt.Errorf("could not create compressed file:", createErr)
	}

	defer createdFile.Close()

	// Compress created file
	compressedFile := CompressToGzipFile(createdFile)

	// Copy original file content
	// to the compressed one
	copyErr := CopyFileContentToCompressedGzipFile(originalFile, compressedFile)

	if copyErr != nil {
		return fmt.Errorf("could not copy original file content to the compressed one:", copyErr)
	}

	defer compressedFile.Close()

	// Delete original file
	deleteErr := DeleteFile(originalFile)

	if deleteErr != nil {
		return fmt.Errorf("Could not delete the original file:", deleteErr)
	}

	return nil
}