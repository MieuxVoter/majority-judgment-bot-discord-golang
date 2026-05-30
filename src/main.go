package main

// A Bot for Discord to create Majority Judgment polls.
// Usage:   /mj create <subject> <proposalA> <proposalB> …

import (
	"flag"
	"fmt"
	"log/slog"
	"main/src/container"
	"main/src/database"
	"main/src/locales"
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
	localizer := locales.GetServerLocalizer()

	// Parse command-line flags
	shouldSyncCommands := flag.Bool( // this one is not overly useful
		"no-sync-commands",
		true,
		"Skip synchronizing commands with Discord",
	)
	flag.Parse()

	// Greet the dev
	fmt.Printf("=== ⚖ MAJORITY JUDGMENT BOT 🤖 v%s ===\n", security.GetVersion())
	//fmt.Printf("Synchronizing commands = %v\n", *shouldSyncCommands)

	// Synchronize the database schema with the Go models
	err := database.Sync()
	if err != nil {
		logger.Error("failed to synchronize db schema", "err", err)
		os.Exit(1)
	}

	// Start the Discord business of the bot
	closeDiscordBot := runner.RunDiscordBot(*shouldSyncCommands)
	defer closeDiscordBot()

	// Perhaps later start the Telegram/Fediverse business of the bot here

	// Finally, start waiting for an interrupting system signal
	logger.Info(localizer.T("FeedbackBotIsRunningPressCtrlC"))
	waitingChannel := make(chan os.Signal, 1)
	signal.Notify(waitingChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-waitingChannel
	logger.Info(localizer.T("FeedbackBotShuttingDown"))
}

func init() {
	// Each service registers into the container in their own init.
	// init() of main is always last, so let's build the container.
	container.Build()

	// Now that the container has been built, let's test our most finicky bits.
	// This helps us catch errors early rather than late.
	runInitChecks()
}

func runInitChecks() {
	logger := container.Get("logger").(*slog.Logger)

	// The rasterizer has a vendored dependency on resvg that can fail in many ways.
	rasterizer := container.Get("rasterizer").(*services.Rasterizer)
	err := rasterizer.Test()
	if err != nil {
		logger.Error("rasterizer fails", "err", err)
	}
}
