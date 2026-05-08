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

func main() {
	// Collect services we're going to use here
	logger := services.GetLogger()

	// Parse command-line flags
	shouldSyncCommands := flag.Bool(
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

	// Read the Discord token from environment
	discordToken, discordTokenFound := os.LookupEnv("DISCORD_TOKEN")
	if !discordTokenFound {
		logger.Fatalln("DISCORD_TOKEN environment variable not found")
	}

	h := handler.New()
	//h.Command("/test", commands.TestHandler)
	//h.Autocomplete("/test", commands.TestAutocompleteHandler)
	h.SlashCommand("/mj", commands.MjDiscordSlashCommandHandler)

	// Connect to Discord and start listening
	client, err := disgo.New(
		discordToken,
		//bot.WithLogger(logger), // requires slog, yet we still use logrus
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
			),
		),
		// For later — http server
		//bot.WithHTTPServerConfigOpts(publicKey,
		//	httpserver.WithAddress(":80"),
		//	httpserver.WithURL("/webhooks/interactions/callback"),
		//),
		bot.WithEventListeners(h),
	)
	if err != nil {
		logger.Fatalln("failed building disgo: %s", err)
	}

	// Close the client's network connection when the bot exits, I guess?
	defer func() {
		closeContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		client.Close(closeContext)
	}()

	// Tell Discord about this bot's available commands, via the REST API
	if *shouldSyncCommands {
		logger.Infoln("Synchronizing commands with Discord…")
		var guilds []snowflake.ID // empty == all guilds
		err = handler.SyncCommands(client, commands.GetDiscordCommands(), guilds)
		if err != nil {
			logger.Fatalln("failed to sync commands: %s", err)
		}
		logger.Infoln("Done synchronizing commands with Discord.")
	}

	gatewayContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.OpenGateway(gatewayContext)
	if err != nil {
		logger.Fatalln(err)
	}

	// Finally, start the waiting loop for a system signal
	logger.Infoln("Bot is running. Press CTRL-C to exit.")
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signalChannel
	logger.Infoln("Shutting down…")
}

func init() {
	// Each service registers into the container in their own init.
	// init() of main is always last, so let's build the container.
	container.Build()
}
