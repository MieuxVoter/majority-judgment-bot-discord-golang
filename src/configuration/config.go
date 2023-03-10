package configuration

// We load our configuration from environment variables.
// They can be defined as usual, or using the .env.local file for convenience.
// It's best if we keep this service clean of dependencies.
// Even logging uses config.  Everything uses config.  Keep this lean.

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"log"
	"main/src/container"
	"os"
)

type Config struct {
	logger *logrus.Logger
}

// Get a configuration value
func (c *Config) Get(key string) string {
	// We could also load from a .ini, if we want?
	value, found := os.LookupEnv(key)
	if !found {
		fmt.Printf("Missing configuration value for `%s'.", key)
	}

	return value
}

// LoadDotEnv loads Environment variables from files, for convenience
func LoadDotEnv() {
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
			LoadDotEnv()
			config := &Config{}

			return config, nil
		},
	})
	if err != nil {
		log.Fatalln("config failed to build", err)
	}
}
