package main

// A Bot for Discord to create polls, using majority judgment.
// Usage:   /mj create <subject> <proposalA> <proposalB> …

import (
	"flag"
	"fmt"
	"main/src/container"
	"main/src/database"
	"main/src/runner"
	"main/src/security"
	"main/src/services"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Collect services we're going to use here
	logger := services.GetLogger()

	// Parse command-line flags
	shouldSyncCommands := flag.Bool( // this one is not overly useful
		"no-sync-commands",
		true,
		"Skip synchronizing commands with Discord",
	)
	flag.Parse()

	// Greet the dev
	fmt.Printf("=== ⚖ MAJORITY JUDGMENT BOT 🤖 v%s ===\n", security.GetVersion())
	fmt.Printf("Synchronizing commands = %v\n", *shouldSyncCommands)

	// Synchronize the database schema with the Go models
	err := database.Sync()
	if err != nil {
		logger.Fatalln(err)
	}

	// Start the Discord business of the bot
	closeDiscordBot := runner.RunDiscordBot(*shouldSyncCommands)
	defer closeDiscordBot()

	// Perhaps later start the Telegram/Fediverse business of the bot here

	// Finally, start waiting for an interrupting system signal
	logger.Infoln("Bot is running. Press CTRL-C to exit.")
	waitingChannel := make(chan os.Signal, 1)
	signal.Notify(waitingChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-waitingChannel
	logger.Infoln("Shutting down…")
}

func init() {
	// Each service registers into the container in their own init.
	// init() of main is always last, so let's build the container.
	container.Build()
}
