package logger

import(
	"strings"
	"os"
	"strconv"
	"fmt"
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
    case "mb":
        return fileSizeInBytes / 1_000_000, nil // 1 million
    case "gb":
        return fileSizeInBytes / 1_000_000_000, nil // 1 billion
    default:
        return fileSizeInBytes, nil
  }
	
}

func GetMaxFileSize() int {
	// Get environment variable
	// for logs file max size
  val := os.Getenv("LOGS_MAX_SIZE")
  if val == "" {
    return 64 // default value
  }

	// Convert string to to type int
  size, err := strconv.Atoi(val)
  if err != nil {
    fmt.Println("Invalid LOGS_MAX_SIZE, using default 64")
    return 64 // default value
  }

  return size
}