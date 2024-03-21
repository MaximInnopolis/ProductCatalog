package logger

import (
	"context"
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

	// Initialize the logger with custom log format
	logger = log.New(io.MultiWriter(logFile, os.Stdout), "", 0) // Remove default date and time prefixes

	// New date and time format with square brackets
	dateTimeFormat := "[2006-01-02 15:04:05] " // Date and time format with square brackets

	// Set a new format for the logger
	logger.SetFlags(0)                                   // Remove default flags
	logger.SetPrefix(dateTimeFormat)                     // Set the prefix for date and time
	logger.SetOutput(io.MultiWriter(logFile, os.Stdout)) // Set the output to file and standard output
}

func Close() error {
	return logFile.Close()
}

func Printf(ctx context.Context, format string, v ...interface{}) {
	endpoint := getEndpointFromRequest(ctx)
	logger.Printf("[%s] "+format, append([]interface{}{endpoint}, v...)...)
}

func Println(v ...interface{}) {
	logger.Println(v...)
}

// getEndpointFromRequest retrieves the endpoint from the context
func getEndpointFromRequest(ctx context.Context) string {
	endpoint, ok := ctx.Value("endpoint").(string)
	if !ok {
		return "unknown_endpoint"
	}
	return endpoint
}
