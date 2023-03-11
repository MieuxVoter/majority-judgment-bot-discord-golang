package services

import (
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"log"
	"main/src/container"
	"os"
)

// GetLogger returns the currently booted logger
func GetLogger() *logrus.Logger {
	return container.Get("logger").(*logrus.Logger)
}

// bootLogger creates a logger for the bot.
// It should be ran AFTER we load .env and .env.local
func bootLogger(config *Config) *logrus.Logger {
	appEnv := config.Get("APP_ENV")
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

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "logger",
		Build: func(ctn di.Container) (interface{}, error) {
			return bootLogger(ctn.Get("config").(*Config)), nil
		},
	})
	if err != nil {
		log.Fatalln("logger failed to build", err)
	}
}
