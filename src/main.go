package main

// A Bot for Discord to create polls, using majority judgment.
// Usage:   /mj create <subject> <proposalA> <proposalB> …

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	cmd "main/src/command"
	db "main/src/database"
	"main/src/logging"
	"os"
)

var log *logrus.Logger

var noCtx = context.Background()

// checkErr logs errors if not nil, along with a user-specified trace
func checkErr(err error, trace string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"trace": trace,
		}).Error(err)
	}
}

// handleMessageMentioningMe reacts when the bot is @ in a channel message
func handleMessageMentioningMe(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	log.Info("Bot has been mentioned: ", msg.Member, msg.Content)

	botUsage := "_I am ready to serve !_   Type `/mj` to start."
	_, err := msg.Reply(noCtx, s, botUsage)
	checkErr(err, "mentioning me")
}

func main() {
	fmt.Println("fgdgfdsgdsgdfs== MAJORITY JUDGMENT BOT v0.0.0 ==") // todo: handle version

	// Load Environment variables from files, for convenience
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("No .env.local file found.  Best create one from .env with your DISCORD_TOKEN.")
	}
	err = godotenv.Load() // .env
	if err != nil {
		fmt.Println("No .env file found.  Ignore this message in builds?")
	}

	// Greet the dev
	fmt.Println("== MAJORITY JUDGMENT BOT v0.0.0 ==") // todo: handle version

	log = logging.BootLogger()

	// Establish a database connection
	_, err = db.Boot(log.Level)
	checkErr(err, "db.Boot")
	err = db.Sync()
	checkErr(err, "db.Sync")

	// Start the Discord client
	client := disgord.New(disgord.Config{
		ProjectName: os.Getenv("DISCORD_NAME"),
		BotToken:    os.Getenv("DISCORD_TOKEN"),
		Logger:      log,
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
	//var permissions disgord.PermissionBit
	u, err := client.BotAuthorizeURL(permissions, []string{
		//"bot", // we're try our best to remove this bot scope, and only use applications.commands
		"applications.commands",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nFollow this URL to authorize the bot on your server:")
	fmt.Println(u)
	fmt.Println("")

	//logFilter, _ := std.NewLogFilter(client)
	filter, _ := std.NewMsgFilter(context.Background(), client)
	//filter.SetPrefix(prefix)

	// Create a handler and bind it to new message events  (no!  use command!)
	//client.Gateway().WithMiddleware(
	//	filter.NotByBot,  // ignore bot messages
	//	filter.HasPrefix, // message must have the given prefix
	//	logFilter.LogMsg,   // log command message
	//	filter.StripPrefix, // remove the command prefix from the message
	//).MessageCreate(handleCommand)

	// Create a handler and bind it to new messages where the bot is mentioned
	// FIXME: check if this still works, and if we can read private direct messages
	client.Gateway().WithMiddleware(
		filter.NotByBot,           // ignore bot messages
		filter.ContainsBotMention, // message must mention this bot
	).MessageCreate(handleMessageMentioningMe)

	// Register slash command once the bot is ready
	//client.Gateway().Ready(func(s disgord.Session, h *disgord.Ready) { // too soon
	client.Gateway().BotReady(func() {
		log.Info("Bot is ready!")
		commands := cmd.GetCommands()
		for _, command := range commands {
			log.Info("Registering command /", command.Name)
			// application command id is 0 here, it's OK.
			// on a ready event, the client is updated to store the application id
			if err = client.ApplicationCommand(0).Global().Create(command); err != nil {
				log.Fatal(err)
			}
		}
	})

	// Respond to discord slash command and other interactions
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		//fmt.Printf("Interaction: %+v\n", *h)
		//fmt.Printf("Data %+v\n", *h.Data)
		//fmt.Printf("Options %+q\n", (*h.Data).Options)

		if h.Type == disgord.InteractionApplicationCommand {

			if len(h.Data.Options) == 0 { // no subcommand was provided
				err = cmd.HandleHelpCommand(noCtx, s, h)
				checkErr(err, "HandleHelpCommand:NoSubcommand")
				return
			}

			// Assumes the subject is the first of our options in the definition, OK if the user changes order
			subCmdName := h.Data.Options[0].Name

			log.Debugln("Handling application command by", h.Member, subCmdName)

			if subCmdName == "help" {
				err = cmd.HandleHelpCommand(noCtx, s, h)
				checkErr(err, "HandleHelpCommand")
			} else if subCmdName == "create" {
				err = cmd.HandleCreateCommand(noCtx, s, h)
				checkErr(err, "HandleCreateCommand")
			} else {
				log.Errorln("Unrecognized subcommand ", subCmdName)
				return
			}

		} else if h.Type == disgord.InteractionMessageComponent {

			if h.Data.ComponentType == disgord.MessageComponentButton {
				log.Debugln("Handling interaction on button", h.Member, h.GuildID)

				var handled = false

				if !handled {
					handled, err = cmd.HandleButtonParticipate(noCtx, s, h)
					checkErr(err, "HandleButtonParticipate")
				}

				if !handled {
					handled, err = cmd.HandleButtonDeliberate(noCtx, s, h)
					checkErr(err, "HandleButtonDeliberate")
				}

				if !handled {
					handled, err = cmd.HandleButtonJudge(noCtx, s, h)
					checkErr(err, "HandleButtonJudge")
				}

				if !handled {
					log.Warnln("Unhandled button interaction", h, h.Data)
					err = cmd.RespondCommandFailure(noCtx, s, h, "This button does nothing.")
					checkErr(err, "RespondCommandFailure:ButtonUnknown")
					return
				}

			} else if h.Data.ComponentType == disgord.MessageComponentSelectMenu {
				// We have no more "select" interactions in the bot for now, but they may come back later.
				// Let's keep this as snippet, it's harmless anyway.
				log.Debugln("Handling interaction on select ", h, h.Data)

				err = s.SendInteractionResponse(context.Background(), h, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackDeferredUpdateMessage,
					Data: &disgord.CreateInteractionResponseData{},
				})
				checkErr(err, "SendInteractionResponse:Select")

			} else {
				log.Warningln("Unhandled interaction on message component ", h, h.Data)
			}

		} else if h.Type == disgord.InteractionPing {
			log.Warningln("Unhandled ping interaction", h, h.Data)
		} else {
			log.Warningln("Unhandled interaction type", h, h.Data)
		}

	})

}
