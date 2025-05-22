package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	// Open original file as
	// reader
	originalFilePath := GetFilePath(originalFile)

	readOriginalPath, readErr := os.Open(originalFilePath)

	if readErr != nil {
		return fmt.Errorf("could not read file %s", originalFilePath)
	}

	// Copy content
	_, copyErr := io.Copy(compressedFile, readOriginalPath)

	fmt.Println("File descriptor:", originalFile.Fd())

	if copyErr != nil {
		return fmt.Errorf("could not copy content: %w", copyErr)
	}

	// Close compressed file
	if err := compressedFile.Close(); err != nil {
		return fmt.Errorf("could not close compressed file: %w", err)
	}

	return nil
}

func DeleteFile(f *os.File) error {
	filePath := GetFilePath(f)

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

func CompressToGzipFile(f *os.File) *gzip.Writer {
	gzWriter := gzip.NewWriter(f)

	return gzWriter
}

func FindFile(root, filename string) (string, error) {
	var foundPath string

	// Search by filename
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
    if err != nil {
      return err
    }
		
    if !d.IsDir() && d.Name() == filename {
      foundPath = path
      return filepath.SkipDir // stop search
    }
    
		return nil
  })

	// Check if error
  if err != nil {
    return "", err
  }

	// Check if path is empty
  if foundPath == "" {
    return "", fmt.Errorf("file %s not found", filename)
  }
  
	return foundPath, nil
}

func CheckIfFileExists(root, filename string) bool {
	_, err := FindFile(root, filename)

	return err != nil
}

func GetFile(root, filename string) (*os.File, error) {
	// Get file path
	filePath, findErr := FindFile(root, filename)

	if findErr != nil {
		return nil, fmt.Errorf("file not found: %w", findErr)
	}

	// Get file
	file, openErr := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if openErr != nil {
		return nil, fmt.Errorf("could not open file: %w", openErr)
	}

	return file, nil
}

func GetFilePath(f *os.File) string {
	return f.Name()
}

func GetFilename(f *os.File) string {
	path := GetFilePath(f)

	pathParts := strings.Split(path, "/")

	// Get last one (filename)
	return pathParts[len(pathParts) - 1]
}



