package logger

import (
	"io"
	"log"
	"os"
)

var (
	logger  *log.Logger
	logFile *os.File
)

const logPath = "logs/app.log"

func init() {
	// Create a directory for log files
	os.MkdirAll("logs", os.ModePerm)
	// Create a file for log messages
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	//defer logFile.Close()

	// Initialize the logger
	logger = log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Close() error {
	return logFile.Close()
}

func Printf(s string, v ...any) {
	logger.Printf(s, v)
}

func Println(v ...any) {
	logger.Println(v)
}
