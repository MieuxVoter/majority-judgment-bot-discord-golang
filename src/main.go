package main

// A Bot for Discord to create polls, using majority judgment.
// Usage:   /mj create <subject> <proposalA> <proposalB> …

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/sirupsen/logrus"
	domain "main/src/command"
	"main/src/container"
	db "main/src/database"
	"main/src/services"
)

var logger *logrus.Logger

var noCtx = context.Background()

// checkErr logs errors if not nil, along with a user-specified trace
func checkErr(err error, trace string) {
	if err != nil {
		logger.WithFields(logrus.Fields{
			"trace": trace,
		}).Error(err)
	}
}

// handleMessageMentioningMe reacts when the bot is @ in a channel message
func handleMessageMentioningMe(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	logger.Info("Bot has been mentioned: ", msg.Member, msg.Content)

	botUsage := "_I am ready to serve !_   Type `/mj` to start."
	_, err := msg.Reply(noCtx, s, botUsage)
	checkErr(err, "mentioning me")
}

func main() {
	// Greet the dev
	fmt.Println("== MAJORITY JUDGMENT BOT v0.0.0 ==") // todo: handle version (govvv?)

	logger = container.Get("logger").(*logrus.Logger)
	config := container.Get("config").(*services.Config)

	// Synchronize the database schema with the Go models
	err := db.Sync()
	checkErr(err, "db.Sync")

	// Start the Discord client
	client := disgord.New(disgord.Config{
		ProjectName: config.Get("DISCORD_NAME"),
		BotToken:    config.Get("DISCORD_TOKEN"),
		Logger:      logger,
		Intents:     disgord.IntentDirectMessages,

		// ! Non-functional due to a current bug, will be fixed upstream someday.
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "Listening to /mj commands",
			},
		},
	})

	// Heartbeat
	defer client.Gateway().StayConnectedUntilInterrupted()

	// Print the link one needs to invite/authorize the bot on their server
	permissions := disgord.PermissionSendMessages |
		disgord.PermissionSendTTSMessages |
		disgord.PermissionSendMessagesInThreads |
		disgord.PermissionAttachFiles |
		disgord.PermissionEmbedLinks
	authorizeURL, err := client.BotAuthorizeURL(permissions, []string{
		//"bot", // we're trying our best to remove this bot scope, and only use applications.commands
		"applications.commands",
	})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("\nFollow this URL to authorize the bot on your server:")
	fmt.Println(authorizeURL)
	fmt.Println("")

	//logFilter, _ := std.NewLogFilter(client)
	filter, _ := std.NewMsgFilter(context.Background(), client)

	// Create a handler and bind it to new (direct) messages
	client.Gateway().WithMiddleware(
		filter.NotByBot, // ignore bot messages
	//	logFilter.LogMsg,   // logger command message
	//	filter.ContainsBotMention, // message must mention this bot
	).MessageCreate(handleMessageMentioningMe)

	// Register slash command once the bot is ready
	//client.Gateway().Ready(func(s disgord.Session, h *disgord.Ready) { // too soon
	client.Gateway().BotReady(func() {
		logger.Info("Bot is ready!")
		commands := domain.GetCommands()
		for _, command := range commands {
			logger.Info("Registering command /", command.Name)
			// application command id is 0 here, it's OK.
			// on a ready event, the client is updated to store the application id
			if err = client.ApplicationCommand(0).Global().Create(command); err != nil {
				logger.Fatal(err)
			}
		}
	})

	// Respond to discord slash command and other interactions
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		//fmt.Printf("Interaction: %+v\n", *h)
		//fmt.Printf("Data %+v\n", *h.Data)
		//fmt.Printf("Options %+q\n", (*h.Data).Options)

		commandInput := domain.DiscordInput{
			Context:     noCtx,
			Session:     s,
			Interaction: h,
		}

		if h.Type == disgord.InteractionApplicationCommand {

			if len(h.Data.Options) == 0 { // no subcommand was provided
				_, err = container.Get("command.help").(*domain.HelpCommand).Handle(commandInput)
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
					commandWasHandled, err = command.Handle(commandInput)
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
					buttonWasHandled, err = button.Handle(commandInput)
					if err != nil {
						checkErr(err, "button handle")
					}
					if buttonWasHandled {
						break
					}
				}

				// TODO: alchemical refactor ; like above, so below
				var handled = buttonWasHandled

				//if !handled {
				//	handled, err = domain.HandleButtonParticipate(noCtx, s, h)
				//	checkErr(err, "HandleButtonParticipate")
				//}

				if !handled {
					handled, err = domain.HandleButtonDeliberate(noCtx, s, h)
					checkErr(err, "HandleButtonDeliberate")
				}

				if !handled {
					handled, err = domain.HandleButtonJudge(noCtx, s, h)
					checkErr(err, "HandleButtonJudge")
				}

				if !handled {
					logger.Warnln("Unhandled button interaction", h, h.Data)
					err = domain.RespondCommandFailure(noCtx, s, h, "This button does nothing.")
					checkErr(err, "RespondCommandFailure:ButtonUnknown")
				}

			} else if h.Data.ComponentType == disgord.MessageComponentSelectMenu {
				// We have no more "select" interactions in the bot for now, but they may come back later.
				// Let's keep this as snippet, it's harmless anyway.
				logger.Debugln("Handling interaction on select ", h, h.Data)

				err = s.SendInteractionResponse(context.Background(), h, &disgord.CreateInteractionResponse{
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

func init() {
	// Each service registers into the container in their own init.
	// init() of main is always last, so let's build the container.
	container.Build()
}
