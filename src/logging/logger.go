package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

// MakeLogger creates a logger for the bot.
// It should be ran AFTER we load .env and .env.local
func MakeLogger() *logrus.Logger {
	appEnv := os.Getenv("APP_ENV")
	logLevel := logrus.DebugLevel
	if appEnv == "prod" {
		logLevel = logrus.InfoLevel
	}

	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logLevel,
	}
}
