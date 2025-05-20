package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func CreateLogFile(folderPath, level string, compressed bool) (*os.File, error){
	// Create folder if not exists
	createFolderErr := os.MkdirAll(folderPath, os.ModePerm)
  if createFolderErr != nil {
    return nil, fmt.Errorf("error creating log folder: %w", createFolderErr)
  }
	
	var filePath string = ""
	levelLower := strings.ToLower(level)

	// Compress and uncompress
	// file's name are different
	if compressed {
		timestamp := time.Now().Format("2006-01-02_15-04-05.000")
		filePath = fmt.Sprintf("%s/%s-compressed-%s.txt.gz", folderPath, levelLower, timestamp)
	} else {
		filePath = fmt.Sprintf("%s/%s.txt", folderPath, levelLower)
	}

	// Create file
	createdFile, createFileErr := os.Create(filePath)
	if createFileErr != nil {
		return nil, fmt.Errorf("error creating log file: %w", createFileErr)
	}
	
	return createdFile, nil
}

func AppendLogToFile(filePath, level, content string) (*os.File, error) {
	// filePath := fmt.Sprintf("%s/%s.txt", folderPath, strings.ToLower(level))
	f, openErr := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

  if openErr != nil {
    return nil, fmt.Errorf("could not open file: %w", openErr)
  }

  defer f.Close()

	_, writeErr := f.WriteString(content)

	if writeErr != nil {
		return nil, fmt.Errorf("could not write in file: %w", writeErr)
	}

	return f, nil
}

func SaveLogToFile(folderPath string, l Logger) error {
	levelLower := strings.ToLower(l.Level.String())
	// Log message format:
	// [Level] Service: Message - (Timestamp)
	logLine := fmt.Sprintf("[%s] %s: %s - (%s)\n", levelLower, l.Service, l.Message, l.Timestamp)
	filePath := fmt.Sprintf("%s/%s.txt", folderPath, levelLower)

	// Find file by path
	foundFile, findErr := GetFile(".", filePath)

	if findErr != nil {
		return fmt.Errorf("file not found: %w", findErr)
	}

	// Get file size in Megabytes
	fileSizeInMegabytes, sizeErr := GetFileSize(foundFile, "mb")

	if sizeErr != nil {
		return fmt.Errorf("could not get file size: %w", sizeErr)
	}

	// Check if file is too big
	if fileSizeInMegabytes >= 64 {
		CompressFile(foundFile, "logs/compressed", levelLower)
	}

	// Append log to file
	file, writeErr := AppendLogToFile(filePath, l.Level, logLine)

	if writeErr != nil {
		return fmt.Errorf("error writing log file: %w", writeErr)
	}

	// Close file
	defer file.Close()

	return nil
}
