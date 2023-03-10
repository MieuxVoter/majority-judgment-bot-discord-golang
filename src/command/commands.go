package command

import (
	"context"
	"github.com/andersfylling/disgord"
	db "main/src/database"
	"main/src/security"
)

// We'll probably want the builder pattern here instead of this static def.  Where is it?
// It would also be nice to figure out how to get variadic commands with unlimited proposals.
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
						Description: "The name of the first proposal, like Friday",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_b",
						Description: "The name of the second proposal, like Pizza",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_c",
						Description: "The name of the third proposal, like Beaujolais",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_d",
						Description: "The name of the fourth proposal, like Michel",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_e",
						Description: "The name of the fifth element, like Moultipass",
						Type:        disgord.OptionTypeString,
					},
					// /!. Discord limits messages integrations to 5 action rows,
					//     so we'd need multiple messages to handle more than 5 proposals.
					//     No point in adding proposal_f here for now, it won't work as-is.
					// > Well, now we use one message per proposal, but how to get variadism here?
				},
			},
			{
				Name:        "help",
				Description: "Send an SOS: ... --- ...",
				Type:        disgord.OptionTypeSubCommand,
			},
		},
	},
}

type Input struct {
	Context     context.Context
	Session     disgord.Session
	Interaction *disgord.InteractionCreate
}

type Command interface {
	Define() *disgord.ApplicationCommandOption
	Matches(command string) bool
	Handle(input *Input) (handled bool, err error)
}

func GetCommands() []*disgord.CreateApplicationCommand {
	return commands
}

// HandleHelpCommand is to refactor into a class
func HandleHelpCommand(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
) error {
	err := s.SendInteractionResponse(ctx, h, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Flags: disgord.MessageFlagEphemeral,
			Content: "🤖 _Hello !_ " +
				"My purpose is to help you create majority judgment polls.\n" +
				"\n" +
				"Try me out:\n" +
				"⌨ `/mj create <subject> <proposal_a> <proposal_b> …`\n" +
				"\n" +
				"⚖ **What is Majority Judgment?**\n" +
				"> A pretty rad **polling system.**  It is used in french :flag_fr: wine 🍷 contests. " +
				"It is simple, subtle and fair.\n" +
				"\n" +
				"🕵 **Can this bot read our messages?**\n" +
				"> **No.**  For extra privacy, this modern bot is NOT allowed to read messages, " +
				"only react to its own `/mj` command and button interactions.\n" +
				"\n" +
				"❺ **Can I use more than 5 grades?**\n" +
				"> **Not for now.**  Discord limits messages to 5 buttons per action row, " +
				"so we'll need more wit to support more grades.\n" +
				"\n" +
				"❺ **Can I use more than 5 proposals?**\n" +
				"> **Not for now.**  Discord does not allow variadic app commands, for now., " +
				"We might have an acceptable workaround, though.\n" +
				"\n" +
				//"\n" +
				//"If \n" +
				"",
		},
	})

	return err
}

// HandleCreateCommand is to refactor into a class at some point
func HandleCreateCommand(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
) error {

	subcommandOptions, err := getSubcommandOptions(h.Data.Options, "create")
	if err != nil {
		return err
	}

	subject := getOptionStringByName(subcommandOptions, "subject", "Poll")
	proposalsNames := make([]string, 0)
	for _, v := range []string{"a", "b", "c", "d", "e"} { // :(|) ooOOk?
		proposalName := getOptionStringByName(subcommandOptions, "proposal_"+v, "")
		if proposalName == "" {
			continue
		}
		proposalsNames = append(proposalsNames, proposalName)
	}

	if len(proposalsNames) < 2 {
		err = RespondCommandFailure(ctx, s, h, "A Poll needs at least two proposals.")
		if err != nil {
			return err
		}
		return nil
	}

	// 8<-----

	guild, err := db.GetOrCreateGuild(db.Engine(), h.GuildID)
	if err != nil {
		return err
	}

	isAllowed, err := security.CanGuildCreatePoll(db.Engine(), guild)
	if err != nil {
		return err
	}
	if !isAllowed {
		err = RespondCommandFailure(ctx, s, h, "This guild cannot create polls anymore.")
		if err != nil {
			return err
		}
		return nil
	}

	// Decrement the guild's quota
	if guild.Quota > 0 {
		guild.Quota = guild.Quota - 1
	}
	_, err = db.Engine().Update(guild)
	if err != nil {
		return err
	}

	poll := &db.Poll{
		Subject: subject,
		GuildId: guild.Id,
	}
	_, err = db.Engine().InsertOne(poll)
	if err != nil {
		return err
	}
	//log.Infoln("New Poll: ", poll.Id, poll.Subject)

	proposals := make([]*db.Proposal, 0)
	for _, proposalName := range proposalsNames {
		proposal := &db.Proposal{
			Name:   proposalName,
			PollId: poll.Id,
		}
		proposals = append(proposals, proposal)
	}
	_, err = db.Engine().Insert(&proposals)
	if err != nil {
		return err
	}

	// 8<-----

	err = RespondWithPollUi(ctx, s, h, poll, proposals, false)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	//fmt.Println("init commands")
}
