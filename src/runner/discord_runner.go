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
	"main/src/locales"
	"main/src/provider"
	"main/src/services"
	"os"
	"time"
)

func RunDiscordBot(
	shouldSyncCommands bool,
) (deferrable func()) {

	// Fetch the services we're going to use
	config := services.GetConfig()
	logger := services.GetLogger()
	localizer := locales.GetServerLocalizer()

	//logHandler := slog.NewTextHandler(
	//	bufio.NewWriter(os.Stdout),
	//	&slog.HandlerOptions{
	//		Level: slog.LevelDebug,
	//	},
	//)
	//slogger := slog.New(logHandler)

	// Read the Discord token from environment
	discordToken := config.Get("DISCORD_TOKEN")
	if discordToken == "" {
		logger.Error("DISCORD_TOKEN environment variable is required.")
		os.Exit(1)
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
				input := provider.DiscordButtonInput{
					Data:  data,
					Event: event,
				}
				handled, err := button.Handle(input)
				if !handled {
					logger.Error("not handled by button", "button", button.GetPattern(), "err", err)
				}
				return err
			},
		)
	}

	// Create the Discord client
	discordClient, err := disgo.New(
		discordToken,
		bot.WithLogger(logger),
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
		logger.Error("failed building disgo", "err", err)
		os.Exit(1)
	}

	// Close the client's network connection when the bot exits
	closeClient := func() {
		closeContext, cancelContext := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelContext()
		discordClient.Close(closeContext)
	}

	// Tell Discord about this bot's available commands, via the REST API
	if shouldSyncCommands {
		logger.Info(localizer.T("FeedbackBotSynchronizingCommands"))
		var guilds []snowflake.ID // empty == all guilds
		err = handler.SyncCommands(discordClient, commands.GetDiscordCommands(), guilds)
		if err != nil {
			logger.Error("failed to sync commands", "err", err)
			os.Exit(1)
		}
		logger.Info(localizer.T("FeedbackBotDoneSynchronizingCommands"))
	}

	// Open a persistent connection to Discord via its Gateway
	gatewayContext, closeGateway := context.WithTimeout(context.Background(), 30*time.Second)
	err = discordClient.OpenGateway(gatewayContext)
	if err != nil {
		logger.Error("failed to open gateway", "err", err)
		os.Exit(1)
	}

	// Define the function the parent must call with defer
	deferrable = func() {
		closeClient()
		closeGateway()
	}

	// Dump the required permissions' integer — dev utility
	//logger.Infoln(
	//	fmt.Sprintf(
	//		"Permissions: %d",
	//		discord.PermissionSendMessages|
	//			discord.PermissionAttachFiles|
	//			discord.PermissionEmbedLinks,
	//		//discord.PermissionSendTTSMessages|
	//		//discord.PermissionSendMessagesInThreads|
	//	),
	//)

	return deferrable

	// OLD STUFF BELOW — can't figure out how to dump the URL using the new disgo

	// Print the link one needs to invite/authorize the bot on their server
	//permissions := disgord.PermissionSendMessages |
	//	disgord.PermissionSendTTSMessages |
	//	disgord.PermissionSendMessagesInThreads |
	//	disgord.PermissionAttachFiles |
	//	disgord.PermissionEmbedLinks
	//authorizeURL, err := discordClient.BotAuthorizeURL(permissions, []string{
	//	"bot",
	//	"applications.commands",
	//})
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//fmt.Println("\nFollow this URL to authorize the bot on your server:")
	//fmt.Println(authorizeURL)
	//fmt.Println("")
}
