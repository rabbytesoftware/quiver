package logger

import (
	"os"
)

func CompressFile(originalFile *os.File, compressFolderPath, level string) error{
	// Create compressed file with name format
	// "level-compressed-timestamp.txt.gz"
	createdFile, createErr := CreateLogFile(compressFolderPath, level, true)

	if createErr != nil {
		return createErr
	}
	
	// Create gzip writer
	compressedFile := CompressToGzipFile(createdFile)
	
	// Copy original file content
	// to the compressed one
	copyErr := CopyFileContentToCompressedGzipFile(originalFile, compressedFile)
	
	if copyErr != nil {
		return copyErr
	}
	
	// Delete original file
	deleteErr := DeleteFile(originalFile)
	
	if deleteErr != nil {
		return deleteErr
	}

	defer createdFile.Close()
	defer compressedFile.Close()

	return nil
}