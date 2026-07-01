// Package logger provides a process-wide structured logger.
package logger

import (
	log "github.com/sirupsen/logrus"
)

// Logger is the global logger instance.
var Logger *log.Entry

// InitLogger initializes the global logger instance.
func InitLogger(loggerInstance *log.Entry) {
	Logger = loggerInstance
}

// LogError logs an error message with context.
func LogError(message string, err error) {
	if Logger != nil {
		Logger.WithError(err).Error(message)
	}
}

// LogErrorf logs a formatted error message with context.
func LogErrorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, args...)
	}
}
