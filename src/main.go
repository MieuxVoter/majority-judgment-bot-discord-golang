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
	// Load Environment variables from files, for convenience
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("No .env.local file found.  Best create one from .env with your DISCORD_TOKEN.")
	}
	err = godotenv.Load() // .env
	if err != nil {
		fmt.Println("No .env file found.  Ignore this message in builds?")
	}

	log = logging.MakeLogger()

	// Greet the dev
	fmt.Println("== MAJORITY JUDGMENT BOT v0.0.0 ==") // todo: handle version

	// Establish a database connection
	_, err = db.Boot(log.Level)
	checkErr(err, "db.Boot")
	err = db.Sync()
	checkErr(err, "db.Sync")

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

	// Note the permission and scope are the minimum requirements for slash command to operate
	u, err := client.BotAuthorizeURL(disgord.PermissionUseSlashCommands, []string{
		"bot", // todo: try our best to remove this scope, and only use application.command
		"applications.command",
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
	client.Gateway().WithMiddleware(
		filter.NotByBot,           // ignore bot messages
		filter.ContainsBotMention, // message must mention this bot
	).MessageCreate(handleMessageMentioningMe)

	// Register slash command once the bot is ready
	client.Gateway().BotReady(func() {
		log.Info("Bot is ready!")
		//log.Info(fmt.Sprintf("Bot %s is ready in application %s"))
		//client.CurrentUser().Get()
		commands := cmd.GetCommands()
		for i := range commands {
			log.Info("Registering command /", commands[i].Name)
			// FIXME: handle multiple guilds; note that config may change slash API
			// - session.GetConnectedGuilds() is empty (?)
			// - needs a database, then
			guildSnow, err := disgord.GetSnowflake("705322981102190593")
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

	// Respond to discord slash command and other interactions
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		//fmt.Printf("Interaction: %+v\n", *h)
		//fmt.Printf("Data %+v\n", *h.Data)
		//fmt.Printf("Options %+q\n", (*h.Data).Options)

		if h.Type == disgord.InteractionApplicationCommand {

			if len(h.Data.Options) == 0 { // no subcommand was provided
				err = s.SendInteractionResponse(context.Background(), h, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Flags: disgord.MessageFlagEphemeral,
						Content: "🤖 _Hello !_  Here are the available command:\n" +
							"⌨ `/mj create <subject> <proposal_a> <proposal_b> …`\n" +
							"⌨ `/mj help`\n" +
							"\n" +
							"> 🕵 For extra privacy; this modern bot is NOT allowed to read messages, " +
							"only react to its own command and interactions.\n" +
							//"\n" +
							//"If \n" +
							"",
					},
				})
				checkErr(err, "SendInteractionResponse:Nothing")
				return
			}

			subCmdName := h.Data.Options[0].Name

			log.Debugln("Handling application command by", h.Member, subCmdName)

			if subCmdName == "help" {
				err = cmd.HandleHelpCommand(noCtx, s, h)
				checkErr(err, "HandleHelpCommand")
				return
			} else if subCmdName == "create" {
				subcommandOptions, err := getSubcommandOptions(h.Data.Options, "create")
				checkErr(err, "getSubcommandOptions:create")

				subject := getOptionStringByName(subcommandOptions, "subject", "Poll")
				proposalsNames := make([]string, 0)
				for _, v := range []string{"a", "b", "c", "d", "e"} { // :(|) ooOOk?
					proposalName := getOptionStringByName(subcommandOptions, "proposal_"+v, "")
					if proposalName == "" {
						continue
					}
					proposalsNames = append(proposalsNames, proposalName)
				}

				if len(proposalsNames) == 0 {
					err = cmd.RespondCommandFailure(noCtx, s, h, "A Poll needs at least two proposals.")
					checkErr(err, "RespondCommandFailure:NoProposals")
					return
				}

				// 8<-----

				poll := db.Poll{
					Subject: subject,
				}
				_, err = db.Orm.InsertOne(&poll)
				checkErr(err, "InsertOne:Poll")
				log.Infoln("New Poll:", poll)

				proposals := make([]*db.Proposal, 0)
				for _, proposalName := range proposalsNames {
					proposal := &db.Proposal{
						Name:   proposalName,
						PollId: poll.Id,
					}
					proposals = append(proposals, proposal)
				}
				_, err = db.Orm.Insert(&proposals)
				checkErr(err, "Insert:Proposals")

				// 8<-----

				pollEmbedHero := &disgord.Embed{
					Title: fmt.Sprintf("⚖ `#%d` %s", poll.Id, subject),
				}
				if len(proposals) > 0 {
					description := ""
					for i, proposal := range proposals {
						if i > 0 {
							description += ", "
						}
						description += proposal.Name
					}
					pollEmbedHero.Description = description
				} else {
					// nothing is cool for now
				}

				err = s.SendInteractionResponse(noCtx, h, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Embeds: []*disgord.Embed{
							pollEmbedHero,
						},
						//Content:    "Bazinga!",
						//SpoilerTagContent:        true,
						// This message might be updated with the merit profile as attachment
						SpoilerTagAllAttachments: true,
						//Components: buttons,
						Components: []*disgord.MessageComponent{
							{
								Type:     disgord.MessageComponentActionRow,
								CustomID: "poll_action_row",
								Components: []*disgord.MessageComponent{
									{
										Type:     disgord.MessageComponentButton,
										Style:    disgord.Success,
										CustomID: fmt.Sprintf("button_participate:%d", poll.Id),
										Label:    "Participate",
										Emoji: &disgord.Emoji{
											Name: "📨",
										},
									},
								},
							},
							//{
							//	Type:     disgord.MessageComponentActionRow,
							//	CustomID: "4",
							//	Components: []*disgord.MessageComponent{
							//		{
							//			Type:        disgord.MessageComponentSelectMenu,
							//			CustomID:    "5",
							//			MinValues:   1,
							//			MaxValues:   1,
							//			Placeholder: "Monday",
							//			Options: []*disgord.SelectMenuOption{
							//				{
							//					Label:       "Excellent",
							//					Description: "I think Monday is Excellent for me",
							//					Value:       "6",
							//					Emoji: &disgord.Emoji{
							//						Name: "🙂",
							//					},
							//				},
							//				{
							//					Label: "Acceptable",
							//					Value: "3",
							//					Emoji: &disgord.Emoji{
							//						Name: "😐",
							//					},
							//				},
							//				{
							//					Label: "Reject",
							//					Value: "0",
							//					Emoji: &disgord.Emoji{
							//						Name: "🙁",
							//					},
							//				},
							//			},
							//		},
							//	},
							//},
						},
					},
				})
				checkErr(err, "SendInteractionResponse:Select")

				return
			} else {
				log.Errorln("Unrecognized subcommand ", subCmdName)

				return
			}

		} else if h.Type == disgord.InteractionMessageComponent {

			if h.Data.ComponentType == disgord.MessageComponentButton {
				log.Debugln("Handling interaction on button", h, h.Data, h.Data.Options)

				var handled = false
				handled, err = cmd.HandleButtonParticipate(noCtx, s, h)
				checkErr(err, "HandleButtonParticipate")

				if !handled {
					handled, err = cmd.HandleButtonJudge(noCtx, s, h)
					checkErr(err, "HandleButtonJudge")
				}

				if !handled {
					log.Warnln("Unhandled button interaction", h, h.Data)
					err = cmd.RespondCommandFailure(noCtx, s, h, "This button does nothing.")
					//err = s.SendInteractionResponse(context.Background(), h, &disgord.CreateInteractionResponse{
					//	Type: disgord.InteractionCallbackChannelMessageWithSource,
					//	Data: &disgord.CreateInteractionResponseData{
					//		Flags:   disgord.MessageFlagEphemeral,
					//		Content: "A Voté. (même pas vrai, c'est en chantier)",
					//	},
					//})
					checkErr(err, "RespondCommandFailure:ButtonUnknown")
					return
				}

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

	})

}
