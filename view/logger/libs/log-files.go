package logger

import (
	"fmt"
	"os"
	"time"
	"compress/gzip"
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

func CreateFile(folderPath, level string, compressed bool) (*os.File, error){
	// Create folder if not exists
	createFolderErr := os.MkdirAll(folderPath, os.ModePerm)
  if createFolderErr != nil {
    return nil, fmt.Errorf("error creating log folder:", createFolderErr)
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
		return nil, fmt.Errorf("error creating log file:", createFileErr)
	}
	
	return createdFile, nil
}

// Do nothing
func AppendLogToFile() (*os.File, error){
	return nil, nil
}

func SaveLogToFile(folderPath string, l Logger){
	// Open file
	file, err := createOrAppendFile(folderPath, l)

	if err != nil {
		fmt.Println("Error opening log file:", err)
		return;
	}

	// Close file
	file.Close()

	// Log message format:
	// [Level] Service: Message - (Timestamp)
  logLine := fmt.Sprintf("[%s] %s: %s - (%s)\n", l.Level, l.Service, l.Message, l.Timestamp)

	// Append log
  _, err = file.WriteString(logLine)

  if err != nil {
    fmt.Println("Error writing to log file:", err)
  }
}
