package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetFileSize(f *os.File, unit string) (int64, error) {
	// Get file info
	fileInfo, err := f.Stat()

	if err != nil {
		return -1, err;
	}

	// Get file size in bytes
	fileSizeInBytes := fileInfo.Size()

	// Improve readability.
	// Working with bytes
	// directly can be messy
	switch strings.ToLower(unit) {
		case "kb":
    	return fileSizeInBytes / 1024, nil
		case "mb":
    	return fileSizeInBytes / (1024 * 1024), nil
		case "gb":
    	return fileSizeInBytes / (1024 * 1024 * 1024), nil
		default:
    	return fileSizeInBytes, nil
	}

}

func GetMaxFileSize() int {
	// Get environment variable
	// for logs file max size
  val := os.Getenv("LOGS_MAX_SIZE")
  if val == "" {
    return 5 // default value
  }

	// Convert string to to type int
  size, err := strconv.Atoi(val)
  if err != nil {
    fmt.Println("Invalid LOGS_MAX_SIZE, using default 64")
    return 64
  }

  return size
}

func IsFileTooLarge(f *os.File) (bool, error) {
	fileSize, err := GetFileSize(f, "kb")

	if err != nil {
		return false, fmt.Errorf("could not get file size: %w", err)
	}

	maxSize := GetMaxFileSize()

	return fileSize >= int64(maxSize), nil
}