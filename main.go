package main

// A Bot for Discord to create polls, using majority judgment.
// Usage:   /mj <subject> <proposalA> <proposalB>

import (
	"context"
	"fmt"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
	//Level:     logrus.InfoLevel,
}

var noCtx = context.Background()

// We'll probably want the builder pattern here instead of this static def
var commands = []*disgord.CreateApplicationCommand{
	{
		Name:        "mj",
		Description: "Manage Majority Judgment polls",
		Type:        disgord.ApplicationCommandChatInput,
		Options: []*disgord.ApplicationCommandOption{
			{
				Name:        "create",
				Description: "Create a new poll",
				Type:        disgord.OptionTypeSubCommand,
				Options: []*disgord.ApplicationCommandOption{
					{
						Name:        "subject",
						Description: "The poll's subject, such as \"When should we meet?\"",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_a",
						Description: "The name of the first proposal",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_b",
						Description: "The name of the second proposal",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_c",
						Description: "The name of the third proposal",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_d",
						Description: "The name of the fourth proposal",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_e",
						Description: "The name of the fifth proposal",
						Type:        disgord.OptionTypeString,
					},
				},
				//Choices: []*disgord.ApplicationCommandOptionChoice{
				//	{
				//		Name:  "create",
				//		Value: "create",
				//	},
				//	{
				//		Name:  "help",
				//		Value: "help",
				//	},
				//},
			},
			{
				Name:        "help",
				Description: "Send an SOS: ... --- ...",
				Type:        disgord.OptionTypeSubCommand,
			},
		},
	},
}

// checkErr logs errors if not nil, along with a user-specified trace
func checkErr(err error, trace string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"trace": trace,
		}).Error(err)
	}
}

// handleCommand is a basic command handler for !commands
// We aim to remove this eventually, and only use application commands (/commands)
func handleCommand(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	//log.Info("Handling Message !")
	//log.Info(msg.Content)
	//log.Info(msg)

	switch msg.Content {
	case "guild":
		_, err := msg.Reply(noCtx, s, msg.GuildID)
		checkErr(err, "guild command")

		log.Info("Connected guilds:")
		for _, guild := range s.GetConnectedGuilds() {
			log.Info("\t" + guild.String())
			_, err := msg.Reply(noCtx, s, guild.String())
			checkErr(err, "guild command")
		}
	case "ping": // whenever the message written is "ping", the bot replies "pong"
		_, err := msg.Reply(noCtx, s, "pong")
		checkErr(err, "ping command")
	case "mj":
		log.Info(msg.Content, msg.Author)

		// Document !poll usage
		pollUsage := fmt.Sprint(
			"_You can create a new poll with this command._\n" +
				"**Usage**: `!mj create <subject>, <proposalA>, <proposalB>, …`\n" +
				"**Example**: `!mj create \"Next meeting, people ?\", Monday, Tuesday, Sunday Morning`")
		_, err := msg.Reply(noCtx, s, pollUsage)
		checkErr(err, "mj command usage")
	default: // unknown command, bot does nothing.
		return
	}
}

// handleMessageMentioningMe reacts when the bot is @ in a channel message
func handleMessageMentioningMe(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	log.Info("Bot has been mentioned: ", msg)

	botUsage := "_I am ready to serve !_   Type `/mj` to start."
	_, err := msg.Reply(noCtx, s, botUsage)
	checkErr(err, "mentioning me")
}

func getSubcommandOptions(
	options []*disgord.ApplicationCommandDataOption,
	name string) ([]*disgord.ApplicationCommandDataOption, error) {

	for _, option := range options {
		if option.Name == name {
			return option.Options, nil
		}
	}

	return nil, fmt.Errorf("command subquery not found")
}

func getOptionStringByName(
	options []*disgord.ApplicationCommandDataOption,
	name string,
	defaultValue string) string {

	for _, option := range options {
		if option.Name == name {
			value := fmt.Sprintf("%s", option.Value)
			if value == "" {
				value = defaultValue
			}

			return value
		}
	}

	return defaultValue
}

func main() {
	const prefix = "!"

	// Load Environment variables from files, for convenience
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Warning("No .env.local file found.  Best create one from .env with your DISCORD_TOKEN.")
	}
	err = godotenv.Load()
	if err != nil {
		log.Warning("No .env file found.  You probably know what you are doing.")
	}

	// Greet the dev
	fmt.Println("== MAJORITY JUDGMENT BOT v0.0.0 ==") // todo: handle version

	// Start the Discord client
	client := disgord.New(disgord.Config{
		ProjectName: "Majority Judgment",
		BotToken:    os.Getenv("DISCORD_TOKEN"),
		Logger:      log,
		Intents: disgord.IntentDirectMessages |
			disgord.IntentGuildMessages |
			disgord.IntentGuildMembers,
		// Remove those once we have what we need
		//disgord.IntentDirectMessageReactions |
		//disgord.IntentDirectMessageTyping |
		//disgord.IntentDirectMessages |
		//disgord.IntentGuildBans |
		//disgord.IntentGuildEmojisAndStickers |
		//disgord.IntentGuildIntegrations |
		//disgord.IntentGuildInvites |
		//disgord.IntentGuildMembers |
		//disgord.IntentGuildMessageReactions |
		//disgord.IntentGuildMessageTyping |
		//disgord.IntentGuildMessages |
		//disgord.IntentGuildPresences |
		//disgord.IntentGuildScheduledEvents |
		//disgord.IntentGuildVoiceStates |
		//disgord.IntentGuildWebhooks |
		//disgord.IntentGuilds |

		// ! Non-functional due to a current bug, will be fixed upstream someday.
		//Presence: &disgord.UpdateStatusPayload{
		//	Game: &disgord.Activity{
		//		Name: "write " + prefix + "ping",
		//	},
		//},
	})

	// Heartbeat
	defer client.Gateway().StayConnectedUntilInterrupted()

	// Note the permission and scope are the minimum requirements for slash commands to operate
	u, err := client.BotAuthorizeURL(disgord.PermissionUseSlashCommands, []string{
		"bot", // todo: try our best to remove this permission, and only use application.commands
		"applications.commands",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Follow this URL to authorize the bot on your server:")
	fmt.Println(u)

	//logFilter, _ := std.NewLogFilter(client)
	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetPrefix(prefix)

	// Create a handler and bind it to new message events
	client.Gateway().WithMiddleware(
		filter.NotByBot,  // ignore bot messages
		filter.HasPrefix, // message must have the given prefix
		//logFilter.LogMsg,   // log command message
		filter.StripPrefix, // remove the command prefix from the message
	).MessageCreate(handleCommand)

	// Create a handler and bind it to new messages where the bot is mentioned
	client.Gateway().WithMiddleware(
		filter.NotByBot,           // ignore bot messages
		filter.ContainsBotMention, // message must mention this bot
	).MessageCreate(handleMessageMentioningMe)

	// Register slash commands once the bot is ready
	client.Gateway().BotReady(func() {
		log.Info("Bot is ready!")
		//log.Info(fmt.Sprintf("Bot %s is ready in aplication %s"))
		//client.CurrentUser().Get()
		for i := range commands {
			log.Info("Registering command /", commands[i].Name)
			// FIXME: handle multiple guilds
			// - session.GetConnectedGuilds() is empty (?)
			// - needs a database, then
			guildSnow, err := disgord.GetSnowflake(705322981102190593)
			checkErr(err, "GetSnowflake:Guild")
			// application command id is 0 here, it's OK.
			// on a ready event, the client is updated to store the application id
			// you can fetch the application id using the bot id (current user id)
			// or copy it from the discord page.
			if err = client.ApplicationCommand(0).Guild(guildSnow).Create(commands[i]); err != nil {
				log.Fatal(err)
			}
		}
	})

	// Respond to discord slash commands and other interactions
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		fmt.Printf("Interaction: %+v\n", *h)
		fmt.Printf("Data %+v\n", *h.Data)
		fmt.Printf("Options %+q\n", (*h.Data).Options)
		//fmt.Printf("Name %s\n", *h.Data.CustomID)
		//fmt.Printf("Options %+v\n", *h.Data.Options)

		if h.Type == disgord.InteractionApplicationCommand {

			log.Debugln("Handling application command ", h, h.Data)

			subcommandOptions, err := getSubcommandOptions(h.Data.Options, "create")
			checkErr(err, "getSubcommandOptions:create")

			subject := getOptionStringByName(subcommandOptions, "subject", "Poll")

			err = s.SendInteractionResponse(context.Background(), h, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					//Content:    "Bazinga!",
					Embeds: []*disgord.Embed{
						{
							Title: "⚖ `#1` " + subject,
							//Description: "Smash that Like button and leave a comment below!",
						},
					},
					//SpoilerTagContent:        true,
					SpoilerTagAllAttachments: true,
					//Components: buttons,
					Components: []*disgord.MessageComponent{
						{
							Type:     disgord.MessageComponentActionRow,
							CustomID: "0",
							Components: []*disgord.MessageComponent{
								{
									Type:     disgord.MessageComponentButton,
									Style:    disgord.Success,
									CustomID: "1",
									Label:    "Good",
									Emoji: &disgord.Emoji{
										Name: "🙂",
									},
								},
								{
									Type:     disgord.MessageComponentButton,
									Style:    disgord.Primary,
									CustomID: "2",
									Label:    "Acceptable",
									Emoji: &disgord.Emoji{
										Name: "😐",
									},
								},
								{
									Type:     disgord.MessageComponentButton,
									Style:    disgord.Danger,
									CustomID: "3",
									Label:    "Reject",
									Emoji: &disgord.Emoji{
										Name: "🙁",
									},
								},
							},
						},
						{
							Type:     disgord.MessageComponentActionRow,
							CustomID: "4",
							Components: []*disgord.MessageComponent{
								{
									Type:        disgord.MessageComponentSelectMenu,
									CustomID:    "5",
									MinValues:   1,
									MaxValues:   1,
									Placeholder: "Monday",
									Options: []*disgord.SelectMenuOption{
										{
											Label:       "Excellent",
											Description: "I think Monday is Excellent for me",
											Value:       "6",
											Emoji: &disgord.Emoji{
												Name: "🙂",
											},
										},
										{
											Label: "Acceptable",
											Value: "3",
											Emoji: &disgord.Emoji{
												Name: "😐",
											},
										},
										{
											Label: "Reject",
											Value: "0",
											Emoji: &disgord.Emoji{
												Name: "🙁",
											},
										},
									},
								},
							},
						},
					},
				},
			})
			checkErr(err, "SendInteractionResponse:Select")

		} else if h.Type == disgord.InteractionMessageComponent {

			if h.Data.ComponentType == disgord.MessageComponentButton {
				log.Debugln("Handling interaction on button", h, h.Data)

				err = s.SendInteractionResponse(context.Background(), h, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Flags:   disgord.MessageFlagEphemeral,
						Content: "A Voté. (même pas vrai, c'est en chantier)",
					},
				})
				checkErr(err, "SendInteractionResponse:Button")

			} else if h.Data.ComponentType == disgord.MessageComponentSelectMenu {
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
			log.Debugln("Unhandled ping interaction", h, h.Data)
		} else {
			log.Warningln("Unhandled interaction type", h, h.Data)
		}

		if h.Data.Type == disgord.ApplicationCommandChatInput {

		}

	})

}
