package logger

import (
	"fmt"
	"os"
)

func SaveLogToFile(folderPath string, l Logger){
	filePath := folderPath + "/" + l.level + ".txt"

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error")
		return;
	}

	
}
