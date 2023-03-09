package command

import (
	"github.com/andersfylling/disgord"
)

import (
	"context"
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

func GetCommands() []*disgord.CreateApplicationCommand {
	return commands
}

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
