package services

import (
	"github.com/sarulabs/di/v2"
	"log"
	"log/slog"
	"main/src/container"
	"os"
)

// GetLogger returns the currently booted logger service.
// This function may only be called *after* init, when the container has been built.
func GetLogger() *slog.Logger {
	return container.Get("logger").(*slog.Logger)
}

// bootLogger creates a logger for the bot.
// It should always be run AFTER we load .env and .env.local, i.e. load the config service.
func bootLogger(config *Config) *slog.Logger {
	appEnv := config.Get("APP_ENV")
	logLevel := slog.LevelDebug
	if appEnv == "prod" {
		logLevel = slog.LevelInfo
	}

	logHandler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: logLevel,
		},
	)

	return slog.New(logHandler)
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
