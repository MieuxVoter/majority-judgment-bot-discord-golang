package main

// A Bot for Discord to create polls, using majority judgment.
// Usage:   /mj create <subject> <proposalA> <proposalB> …

import (
	"context"
	"flag"
	"fmt"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"main/src/commands"
	"main/src/container"
	"main/src/database"
	"main/src/security"
	"main/src/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger *logrus.Logger

func main() {
	// Collect services we're going to use here
	logger = services.GetLogger()

	// Parse command-line flags
	shouldSyncCommands := flag.Bool(
		"sync-commands",
		false,
		"Whether to synchronize commands with discord",
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

	// Read the Discord token from environment
	discordToken, discordTokenFound := os.LookupEnv("DISCORD_TOKEN")
	if !discordTokenFound {
		logger.Fatalln("DISCORD_TOKEN environment variable not found")
	}

	h := handler.New()
	h.Command("/test", commands.TestHandler)
	//h.Autocomplete("/test", commands.TestAutocompleteHandler)

	// Connect to Discord and start listening
	client, err := disgo.New(
		discordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
			),
		),
		//bot.WithEventListenerFunc(func(e *events.MessageCreate) {
		//	// event code here
		//}),
		bot.WithEventListeners(h),
	)
	if err != nil {
		logger.Fatalln(err)
	}

	// Close the client's network connection when the bot exits, I guess?
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		client.Close(ctx)
	}()

	if *shouldSyncCommands {
		logger.Infoln("Synchronizing commands with Discord…")
		//slog.Info("Syncing commands", slog.Any("guild_ids", cfg.Bot.DevGuilds))
		guilds := []snowflake.ID{} // empty == globally
		err = handler.SyncCommands(client, commands.Commands, guilds)
		if err != nil {
			logger.Fatalln("Failed to sync commands", err)
		}
		logger.Infoln("Done synchronizing commands with Discord.")
	}

	gatewayContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = client.OpenGateway(gatewayContext); err != nil {
		logger.Fatalln(err)
	}

	// Finally, start the waiting loop
	logger.Infoln("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	logger.Infoln("Shutting down…")
}

func init() {
	// Each service registers into the container in their own init.
	// init() of main is always last, so let's build the container.
	container.Build()
}
