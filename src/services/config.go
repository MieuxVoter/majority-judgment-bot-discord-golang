package services

// We load our configuration from environment variables.
// They can be defined as usual, or using the .env.local file for convenience.

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	"os"
)

// Config is the service handling bot configuration retrieval.
// It's best if we keep this service clean of dependencies.
// Even logging uses config.  Everything uses config.  Keep this lean.
type Config struct{}

// Get a configuration value.
func (c *Config) Get(key string) string {
	value, found := os.LookupEnv(key)
	if !found {
		fmt.Printf("WARNING: Missing configuration value for `%s'.\n", key)
	}
	// We could also load from a .ini, here, if we wanted to.

	return value
}

// GetConfig returns the configuration service.
// This function may only be called *after* init, when the container has been built.
func GetConfig() *Config {
	return container.Get("config").(*Config)
}

// loadDotEnv loads Environment variables from files, for convenience.
func loadDotEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("No .env.local file found.  Best create one from .env with your DISCORD_TOKEN.")
	}
	err = godotenv.Load() // .env
	if err != nil {
		fmt.Println("No .env file found.  Ignore this message in builds?")
	}
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "config",
		Build: func(container di.Container) (interface{}, error) {
			loadDotEnv()
			config := &Config{}

			return config, nil
		},
	})
	if err != nil {
		log.Fatalln("config failed to build", err)
	}
}
