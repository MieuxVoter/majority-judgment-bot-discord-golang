package runner

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"main/src/commands"
	"main/src/container"
	"main/src/domain"
	"main/src/provider"
	"main/src/services"
	"time"
)

// Wo do not need these memoization(s).  Remove these comments at some point.

// config is a local memoization of the [services.Config] service.
// It requires that the function [RunDiscordBot] is run first.
//var config *services.Config

// logger is a local memoization of the [logrus.Logger] service.
// It requires that the function [RunDiscordBot] is run first.
//var logger *logrus.Logger

// discordClient connects to the REST HTTP API of Discord.
// It requires that the function [RunDiscordBot] is run first.
//var discordClient *bot.Client

// checkErr logs errors if not nil, along with a user-specified trace
//func checkErr(err error, trace string) {
//	if err != nil {
//		logger.WithFields(logrus.Fields{
//			"trace": trace,
//		}).Error(err)
//	}
//}

func RunDiscordBot(
	shouldSyncCommands bool,
) (deferrable func()) {

	// Fetch the services we're going to use
	config := services.GetConfig()
	logger := services.GetLogger()

	// Read the Discord token from environment
	discordToken := config.Get("DISCORD_TOKEN")
	if discordToken == "" {
		logger.Fatalln("DISCORD_TOKEN environment variable not found")
	}

	// Register our slash command(s)
	h := handler.New()
	h.SlashCommand("/mj", commands.MjDiscordSlashCommandHandler)

	// Register our button handlers
	for _, buttonGeneric := range container.GetCollection("button.") {
		button := buttonGeneric.(domain.Button)
		h.ButtonComponent(
			button.GetPattern(),
			func(
				data discord.ButtonInteractionData,
				event *handler.ComponentEvent,
			) error {
				//logger.Debugln("button handled:", button.GetPattern())
				input := provider.DiscordButtonInput{
					Data:  data,
					Event: event,
				}
				handled, err := button.Handle(input)
				if !handled {
					logger.Errorln("not handled by button", button.GetPattern(), err)
				}
				return err
			},
		)
	}

	// Create the Discord client
	discordClient, err := disgo.New(
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
		logger.Fatalln("failed building disgo:", err)
	}

	// Close the client's network connection when the bot exits
	closeClient := func() {
		closeContext, cancelContext := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelContext()
		discordClient.Close(closeContext)
	}

	// Tell Discord about this bot's available commands, via the REST API
	if shouldSyncCommands {
		logger.Infoln("Synchronizing commands with Discord…")
		var guilds []snowflake.ID // empty == all guilds
		err := handler.SyncCommands(discordClient, commands.GetDiscordCommands(), guilds)
		if err != nil {
			logger.Fatalln("failed to sync commands: %s", err)
		}
		logger.Infoln("Done synchronizing commands with Discord.")
	}

	// Open a persistent connection to Discord via its Gateway
	gatewayContext, closeGateway := context.WithTimeout(context.Background(), 30*time.Second)
	err = discordClient.OpenGateway(gatewayContext)
	if err != nil {
		logger.Fatalln("failed to open gateway:", err)
	}

	// Define the function the parent must call with defer
	deferrable = func() {
		closeClient()
		closeGateway()
	}

	return deferrable

	// OLD STUFF BELOW

	// Print the link one needs to invite/authorize the bot on their server
	//permissions := disgord.PermissionSendMessages |
	//	disgord.PermissionSendTTSMessages |
	//	disgord.PermissionSendMessagesInThreads |
	//	disgord.PermissionAttachFiles |
	//	disgord.PermissionEmbedLinks
	//authorizeURL, err := discordClient.BotAuthorizeURL(permissions, []string{
	//	//"bot", // we're trying our best to remove this bot scope, and only use applications.commands
	//	"applications.commands",
	//})
	//if err != nil {
	//	logger.Fatal(err)
	//}

	//fmt.Println("\nFollow this URL to authorize the bot on your server:")
	//fmt.Println(authorizeURL)
	//fmt.Println("")
}
