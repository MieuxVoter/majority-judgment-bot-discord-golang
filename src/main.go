package main

// A Bot for Discord to create polls, using majority judgment.
// Usage:   /mj create <subject> <proposalA> <proposalB> …

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"main/src/container"
	"main/src/database"
	"main/src/provider/discord"
	"main/src/security"
	"main/src/services"
)

var logger *logrus.Logger

func main() {
	// Collect services we're going to use
	logger = services.GetLogger()

	// Greet the dev
	fmt.Printf("=== ⚖  MAJORITY JUDGMENT BOT 🤖 v%s ===\n", security.GetVersion())

	// Synchronize the database schema with the Go models
	err := database.Sync()
	if err != nil {
		logger.Fatalln(err)
	}

	// Connect to Discord and start listening
	discord.Run()
}

func init() {
	// Each service registers into the container in their own init.
	// init() of main is always last, so let's build the container.
	container.Build()
}
