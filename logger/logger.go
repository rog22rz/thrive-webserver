package logger

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
)

var (
	client *logging.Client
	logger *logging.Logger
)

// Init initializes the logger.
func Init(projectID, logName string) error {
	ctx := context.Background()
	var err error
	client, err = logging.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	logger = client.Logger(logName)
	return nil
}

// Log logs a message.
func Log(message string) {
	if logger != nil {
		logger.Log(logging.Entry{Payload: message})
	} else {
		log.Printf("Logger is not initialized: %s", message)
	}
}

// Close closes the logger client.
func Close() {
	if client != nil {
		client.Close()
	}
}
