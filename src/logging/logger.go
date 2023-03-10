package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger = nil

// BootLogger creates a logger for the bot.
// It should be ran AFTER we load .env and .env.local
func BootLogger() *logrus.Logger {
	if logger != nil {
		return logger // idempotent
	}

	appEnv := os.Getenv("APP_ENV")
	logLevel := logrus.DebugLevel
	if appEnv == "prod" {
		logLevel = logrus.InfoLevel
	}

	logger = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logLevel,
	}

	return logger
}

// GetLogger returns the currently booted logger
func GetLogger() *logrus.Logger {
	return logger
}
