package runner

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"main/src/commands"
	"main/src/container"
	"main/src/domain"
	"main/src/provider"
	"main/src/services"
	"time"
)

// noCtx is probably a horrible, noob design oversight on my part. Bleh.
//var noCtx = context.Background()

// config is a local memoization of the [services.Config] service.
// It requires that the function [RunDiscordBot] is run first. (quite bad design)
var config *services.Config

// logger is a local memoization of the [logrus.Logger] service.
// It requires that the function [RunDiscordBot] is run first. (quite bad design)
var logger *logrus.Logger

var discordClient *bot.Client

// checkErr logs errors if not nil, along with a user-specified trace
func checkErr(err error, trace string) {
	if err != nil {
		logger.WithFields(logrus.Fields{
			"trace": trace,
		}).Error(err)
	}
}

func registerDiscordCommands(client *bot.Client) {
	logger.Infoln("Synchronizing commands with Discord…")
	var guilds []snowflake.ID // empty == all guilds
	err := handler.SyncCommands(client, commands.GetDiscordCommands(), guilds)
	if err != nil {
		logger.Fatalln("failed to sync commands: %s", err)
	}
	logger.Infoln("Done synchronizing commands with Discord.")
}

//func onBotReady() {
//	logger.Info("Bot is readying…")
//	registerDiscordCommands(discordClient)
//	logger.Info("Bot is ready.")
//}

//func onInteraction(s disgord.Session, h *disgord.InteractionCreate) {
//	// Handy debug/exploration snippets
//	//fmt.Printf("Interaction: %+v\n", *h)
//	//fmt.Printf("Data %+v\n", *h.Data)
//	//fmt.Printf("Options %+q\n", (*h.Data).Options)
//	//fmt.Printf("%+q\n", h.GuildID)
//	//fmt.Printf("%+q\n", h.ChannelID)
//
//	var err error
//	vendorInput := provider.DiscordCommandInput{
//		Context:     noCtx,
//		Session:     s,
//		Interaction: h,
//	}
//
//	if h.Type == disgord.InteractionApplicationCommand {
//
//		if len(h.Data.Options) == 0 { // no subcommand was provided
//			_, err = container.Get("command.help").(*domain.HelpCommand).Handle(vendorInput)
//			checkErr(err, "HandleHelpCommand:NoSubcommand")
//			return
//		}
//
//		subCmdName := h.Data.Options[0].Name
//
//		logger.Debugln("Handling application command by", h.Member, subCmdName)
//
//		commands := container.GetCollection("command")
//		commandWasHandled := false
//		for _, commandGeneric := range commands {
//			command := commandGeneric.(domain.Command)
//			if command.Matches(subCmdName) {
//				commandWasHandled, err = command.Handle(vendorInput)
//				if err != nil {
//					checkErr(err, "command "+subCmdName)
//				}
//				if commandWasHandled {
//					break
//				}
//			}
//		}
//
//		if !commandWasHandled {
//			logger.Errorln("Unrecognized subcommand ", subCmdName)
//		}
//
//	} else if h.Type == disgord.InteractionMessageComponent {
//
//		if h.Data.ComponentType == disgord.MessageComponentButton {
//
//			logger.Debugln("Handling interaction on button", h.Member, h.GuildID)
//
//			buttons := container.GetCollection("button")
//			buttonWasHandled := false
//			for _, buttonGeneric := range buttons {
//				button := buttonGeneric.(domain.Button)
//				buttonWasHandled, err = button.Handle(vendorInput)
//				if err != nil {
//					checkErr(err, "button handle")
//				}
//				if buttonWasHandled {
//					break
//				}
//			}
//
//			if !buttonWasHandled {
//				logger.Warnln("Unhandled button interaction", h, h.Data)
//				err = domain.RespondServerError(vendorInput, "This button does nothing.")
//				checkErr(err, "RespondCommandFailure:ButtonUnknown")
//			}
//
//		} else if h.Data.ComponentType == disgord.MessageComponentSelectMenu {
//			// We have no more "select" interactions in the bot for now, but they may come back later.
//			// Let's keep this as snippet, it's harmless anyway.
//			logger.Debugln("Handling interaction on select ", h, h.Data)
//
//			err = s.SendInteractionResponse(noCtx, h, &disgord.CreateInteractionResponse{
//				Type: disgord.InteractionCallbackDeferredUpdateMessage,
//				Data: &disgord.CreateInteractionResponseData{},
//			})
//			checkErr(err, "SendInteractionResponse:Select")
//
//		} else {
//			logger.Warningln("Unhandled interaction on message component ", h, h.Data)
//		}
//
//	} else if h.Type == disgord.InteractionPing {
//		logger.Warningln("Unhandled ping interaction", h, h.Data)
//	} else {
//		logger.Warningln("Unhandled interaction type", h, h.Data)
//	}
//}

// onDirectMessageToMe reacts when the bot is @ in a channel message
//func onDirectMessageToMe(s disgord.Session, data *disgord.MessageCreate) {
//	msg := data.Message
//
//	logger.Info("Bot has been talked to: ", msg.Member, msg.Content)
//
//	botUsage := "🤖 _I am ready to assist._   Type `/mj` to start."
//	_, err := msg.Reply(noCtx, s, botUsage)
//	checkErr(err, "onDirectMessageToMe")
//}

func RunDiscordBot(
	shouldSyncCommands bool,
) (deferrable func()) {

	// Fetch and memoize services we're going to use.
	config = services.GetConfig()
	logger = services.GetLogger()

	// Read the Discord token from environment
	discordToken := config.Get("DISCORD_TOKEN")
	if discordToken == "" {
		logger.Fatalln("DISCORD_TOKEN environment variable not found")
	}

	h := handler.New()
	h.SlashCommand("/mj", commands.MjDiscordSlashCommandHandler)

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
					logger.Errorln("not handled by button", button.GetPattern())
				}
				return err
			},
		)
	}

	//h.ButtonComponent("/button/poll/vote/{id}", func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	//	logger.Infoln("button handled!!")
	//	return nil
	//})
	//h.Command("/test", commands.TestHandler)
	//h.Autocomplete("/test", commands.TestAutocompleteHandler)

	// Connect to Discord and start listening
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
		logger.Fatalln("failed building disgo: %s", err)
	}

	// Close the client's network connection when the bot exits.
	cancelClient := func() {
		closeContext, cancelContext := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelContext()
		discordClient.Close(closeContext)
	}

	// Tell Discord about this bot's available commands, via the REST API.
	if shouldSyncCommands {
		registerDiscordCommands(discordClient)
	}

	gatewayContext, cancelGateway := context.WithTimeout(context.Background(), 30*time.Second)
	err = discordClient.OpenGateway(gatewayContext)
	if err != nil {
		logger.Fatalln("failed to open gateway:", err)
	}

	deferrable = func() {
		cancelClient()
		cancelGateway()
	}

	return deferrable

	// OLD STUFF BELOW

	// Start the Discord client
	//discordClient = disgord.New(disgord.Config{
	//	ProjectName: config.Get("DISCORD_NAME"),
	//	BotToken:    config.Get("DISCORD_TOKEN"),
	//	Logger:      logger,
	//	Intents:     disgord.IntentDirectMessages,
	//
	//	// ! Non-functional due to a current bug in disgord, will be fixed upstream someday.
	//	Presence: &disgord.UpdateStatusPayload{
	//		Game: &disgord.Activity{
	//			Name: "Ranking proposals (`/mj`)",
	//		},
	//	},
	//})

	// Heartbeat
	//defer discordClient.Gateway().StayConnectedUntilInterrupted()

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

	//filter, _ := std.NewMsgFilter(noCtx, discordClient)
	//logFilter, _ := std.NewLogFilter(discordClient)

	// Register slash command once the bot is ready
	//discordClient.Gateway().BotReady(onBotReady)

	// Handle slash command and other interactions (buttons, selects, …)
	//discordClient.Gateway().InteractionCreate(onInteraction)

	// Handle (direct) messages
	//discordClient.Gateway().WithMiddleware(
	//	filter.NotByBot, // ignore bot messages
	//	//logFilter.LogMsg,  // logger command message
	//).MessageCreate(onDirectMessageToMe)
}
