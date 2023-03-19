package discord

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/sirupsen/logrus"
	"main/src/container"
	"main/src/database"
	"main/src/domain"
	"main/src/security"
	"main/src/services"
)

var logger *logrus.Logger
var noCtx = context.Background()
var discordClient *disgord.Client

// checkErr logs errors if not nil, along with a user-specified trace
func checkErr(err error, trace string) {
	if err != nil {
		logger.WithFields(logrus.Fields{
			"trace": trace,
		}).Error(err)
	}
}

// handleDirectMessageToMe reacts when the bot is @ in a channel message
func handleDirectMessageToMe(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	logger.Info("Bot has been talked to: ", msg.Member, msg.Content)

	botUsage := "🤖 _I am ready to assist._   Type `/mj` to start."
	_, err := msg.Reply(noCtx, s, botUsage)
	checkErr(err, "mentioning me")
}

func onDiscordBotReady() {
	logger.Info("Bot is ready!")
	commands := domain.GetCommands()
	for _, command := range commands {
		logger.Info("Registering command /", command.Name)
		// application command id is 0 here, it's OK.
		// on a ready event, the client is updated to store the application id
		if err := discordClient.ApplicationCommand(0).Global().Create(command); err != nil {
			logger.Fatal(err)
		}
	}
}

func Run() {
	// Greet the dev
	fmt.Printf("=== ⚖  MAJORITY JUDGMENT BOT 🤖 v%s ===\n", security.GetVersion())

	logger = container.Get("logger").(*logrus.Logger)
	config := container.Get("config").(*services.Config)

	// Synchronize the database schema with the Go models
	err := database.Sync()
	checkErr(err, "db.Sync")

	// Start the Discord client
	discordClient = disgord.New(disgord.Config{
		ProjectName: config.Get("DISCORD_NAME"),
		BotToken:    config.Get("DISCORD_TOKEN"),
		Logger:      logger,
		Intents:     disgord.IntentDirectMessages,

		// ! Non-functional due to a current bug, will be fixed upstream someday.
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "Ranking proposals (`/mj`)",
			},
		},
	})

	// Heartbeat
	defer discordClient.Gateway().StayConnectedUntilInterrupted()

	// Print the link one needs to invite/authorize the bot on their server
	permissions := disgord.PermissionSendMessages |
		disgord.PermissionSendTTSMessages |
		disgord.PermissionSendMessagesInThreads |
		disgord.PermissionAttachFiles |
		disgord.PermissionEmbedLinks
	authorizeURL, err := discordClient.BotAuthorizeURL(permissions, []string{
		//"bot", // we're trying our best to remove this bot scope, and only use applications.commands
		"applications.commands",
	})
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println("\nFollow this URL to authorize the bot on your server:")
	fmt.Println(authorizeURL)
	fmt.Println("")

	filter, _ := std.NewMsgFilter(noCtx, discordClient)
	//logFilter, _ := std.NewLogFilter(discordClient)

	// Create a handler and bind it to new (direct) messages
	discordClient.Gateway().WithMiddleware(
		filter.NotByBot, // ignore bot messages
		//	logFilter.LogMsg,   // logger command message
	).MessageCreate(handleDirectMessageToMe)

	// Register slash command once the bot is ready
	discordClient.Gateway().BotReady(onDiscordBotReady)

	// Respond to discord slash command and other interactions
	discordClient.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		// Handy debug/exploration snippets
		//fmt.Printf("Interaction: %+v\n", *h)
		//fmt.Printf("Data %+v\n", *h.Data)
		//fmt.Printf("Options %+q\n", (*h.Data).Options)
		//fmt.Printf("%+q\n", h.GuildID)
		//fmt.Printf("%+q\n", h.ChannelID)

		vendorInput := domain.DiscordInput{
			Context:     noCtx,
			Session:     s,
			Interaction: h,
		}

		if h.Type == disgord.InteractionApplicationCommand {

			if len(h.Data.Options) == 0 { // no subcommand was provided
				_, err = container.Get("command.help").(*domain.HelpCommand).Handle(vendorInput)
				if err != nil {
					checkErr(err, "HandleHelpCommand:NoSubcommand")
				}
				return
			}

			subCmdName := h.Data.Options[0].Name

			logger.Debugln("Handling application command by", h.Member, subCmdName)

			commands := container.GetCollection("command")
			commandWasHandled := false
			for _, commandGeneric := range commands {
				command := commandGeneric.(domain.Command)
				if command.Matches(subCmdName) {
					commandWasHandled, err = command.Handle(vendorInput)
					if err != nil {
						checkErr(err, "command "+subCmdName)
					}
					if commandWasHandled {
						break
					}
				}
			}

			if !commandWasHandled {
				logger.Errorln("Unrecognized subcommand ", subCmdName)
			}

		} else if h.Type == disgord.InteractionMessageComponent {

			if h.Data.ComponentType == disgord.MessageComponentButton {

				logger.Debugln("Handling interaction on button", h.Member, h.GuildID)

				buttons := container.GetCollection("button")
				buttonWasHandled := false
				for _, buttonGeneric := range buttons {
					button := buttonGeneric.(domain.Button)
					buttonWasHandled, err = button.Handle(vendorInput)
					if err != nil {
						checkErr(err, "button handle")
					}
					if buttonWasHandled {
						break
					}
				}

				if !buttonWasHandled {
					logger.Warnln("Unhandled button interaction", h, h.Data)
					err = domain.RespondServerError(vendorInput, "This button does nothing.")
					checkErr(err, "RespondCommandFailure:ButtonUnknown")
				}

			} else if h.Data.ComponentType == disgord.MessageComponentSelectMenu {
				// We have no more "select" interactions in the bot for now, but they may come back later.
				// Let's keep this as snippet, it's harmless anyway.
				logger.Debugln("Handling interaction on select ", h, h.Data)

				err = s.SendInteractionResponse(noCtx, h, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackDeferredUpdateMessage,
					Data: &disgord.CreateInteractionResponseData{},
				})
				checkErr(err, "SendInteractionResponse:Select")

			} else {
				logger.Warningln("Unhandled interaction on message component ", h, h.Data)
			}

		} else if h.Type == disgord.InteractionPing {
			logger.Warningln("Unhandled ping interaction", h, h.Data)
		} else {
			logger.Warningln("Unhandled interaction type", h, h.Data)
		}

	})

}
