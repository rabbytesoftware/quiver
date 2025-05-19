package logger

import (
	"fmt"
	"os"
	"time"
)

func createOrAppendFile(folderPath string, l Logger) (*os.File, error) {
	
	err := os.MkdirAll(folderPath, os.ModePerm)
  if err != nil {
    fmt.Println("Error creating log folder:", err)
    return nil, err
  }
	
	filePath := folderPath + "/" + l.Level + ".txt"
	
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	return file, nil;
}

func CreateLogFile(folderPath, level string, compressed bool) (*os.File, error){
	// Create folder if not exists
	createFolderErr := os.MkdirAll(folderPath, os.ModePerm)
  if createFolderErr != nil {
    return nil, fmt.Errorf("error creating log folder: %w", createFolderErr)
  }
	
	var filePath string = ""

	// Compress and uncompress
	// file's name are different
	if compressed {
		timestamp := time.Now().Format("2006-01-02_15-04-05.000")
		filePath = fmt.Sprintf("%s/%s-compressed-%s.txt.gz", folderPath, level, timestamp)
	} else {
		filePath = fmt.Sprintf("%s/%s.txt", folderPath, level)
	}

	// Create file
	createdFile, createFileErr := os.Create(filePath)
	if createFileErr != nil {
		return nil, fmt.Errorf("error creating log file: %w", createFileErr)
	}
	
	return createdFile, nil
}

func AppendLogToFile(folderPath, level, content string) error {
	filePath := fmt.Sprintf("%s/%s.txt", folderPath, level)
	f, openErr := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

  if openErr != nil {
    return fmt.Errorf("could not open file: %w", openErr)
  }

  defer f.Close()

	_, writeErr := f.WriteString(content)

	if writeErr != nil {
		return fmt.Errorf("could not write in file: %w", writeErr)
	}

	return nil
}

func SaveLogToFile(folderPath string, l Logger) error {
	// Open file
	file, openErr := createOrAppendFile(folderPath, l)

	if openErr != nil {
		return fmt.Errorf("Error opening log file: %w", openErr)
	}

	// Close file
	defer file.Close()

	// Log message format:
	// [Level] Service: Message - (Timestamp)
  logLine := fmt.Sprintf("[%s] %s: %s - (%s)\n", l.Level, l.Service, l.Message, l.Timestamp)

	// Append log
  _, writeErr := file.WriteString(logLine)

  if writeErr != nil {
		return fmt.Errorf("error writing to log file: %w", writeErr)
  }

	return nil
}
