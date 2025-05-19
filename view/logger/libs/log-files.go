package logger

import (
	"fmt"
	"os"
)

func createOrAppendFile(folderPath, l Logger) (*os.File, error) {
	
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
