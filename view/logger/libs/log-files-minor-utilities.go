package logger

import (
	"io"
	"os"
	"fmt"
	"compress/gzip"
)

func CopyFileContent(originalFile, newFile *os.File) error {
	// Move the read pointer to the
	// beginning of the original file
	_, seekErr := originalFile.Seek(0, 0)
	if seekErr != nil {
		return fmt.Errorf("could not seek file: %w", seekErr)
	}

	// Copy content
	_, copyErr := io.Copy(newFile, originalFile)

	if copyErr != nil {
		return fmt.Errorf("could not copy content: %w", copyErr)
	}

	return nil
}

func CopyFileContentToCompressedGzipFile(originalFile *os.File, compressedFile *gzip.Writer) error {
	// Move the read pointer to the
	// beginning of the original file
	_, seekErr := originalFile.Seek(0, 0)
	if seekErr != nil {
		return fmt.Errorf("could not seek file: %w", seekErr)
	}

	// Copy content
	_, copyErr := io.Copy(compressedFile, originalFile)

	if copyErr != nil {
		return fmt.Errorf("could not copy content: %w", copyErr)
	}

	return nil
}

func DeleteFile(f *os.File) error {
	filePath := f.Name()

	// Attempt to close the file before deleting it
	if closeErr := f.Close(); closeErr != nil {
		return fmt.Errorf("failed to close file: %w", closeErr)
	}

	// Attempt to delete the file after closing
	if deleteErr := os.Remove(filePath); deleteErr != nil {
		return fmt.Errorf("failed to delete file: %w", deleteErr)
	}

	return nil
}

func CompressToGzipFile(f *os.File) *gzip.Writer{
	gzWriter := gzip.NewWriter(f)

	return gzWriter
}