package command

import (
	"fmt"
	"github.com/andersfylling/disgord"
	db "main/src/database"
)

import (
	"context"
)

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
				"Here are the available command:\n" +
				"⌨ `/mj create <subject> <proposal_a> <proposal_b> …`\n" +
				"⌨ `/mj help`\n" +
				"\n" +
				"🕵 **Can this bot read our messages?**\n" +
				"> **No.**  For extra privacy, this modern bot is NOT allowed to read messages, " +
				"only react to its own `/mj` command and button interactions.\n" +
				"\n" +
				"❺ **Can I use more than 5 proposals?**\n" +
				"> **Not for now.**  Discord limits messages to 5 action rows, " +
				"so we'll need more code to support more proposals.\n" +
				"\n" +
				//"\n" +
				//"If \n" +
				"",
		},
	})

	return err
}

func StartVoteProcess(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
	judge *disgord.Member,
	proposal *db.Proposal,
	poll *db.Poll,
	previousJudgment *db.Judgment,
) error {

	interactionResponse := &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Flags: disgord.MessageFlagEphemeral,
			Embeds: []*disgord.Embed{
				{
					Title:       fmt.Sprintf("⚖ `#%d` %s", poll.Id, proposal.Name),
					Description: fmt.Sprintf("What do you think of **_%s_** ?", proposal.Name),
				},
			},
			Components: []*disgord.MessageComponent{
				{
					Type:       disgord.MessageComponentActionRow,
					CustomID:   "poll_action_row",
					Components: []*disgord.MessageComponent{}, // filled below
				},
			},
		},
	}

	for gradeLevel, grade := range poll.GetGradingSlice() {

		previouslySelectedMarker := ""
		if previousJudgment != nil {
			if uint8(gradeLevel) == previousJudgment.Grade {
				previouslySelectedMarker = " ✅"
			}
		}
		interactionResponse.Data.Components[0].Components = append(
			interactionResponse.Data.Components[0].Components,
			&disgord.MessageComponent{
				Type:     disgord.MessageComponentButton,
				Style:    disgord.Primary,
				CustomID: fmt.Sprintf("button_judge:%d:%d", proposal.Id, gradeLevel),
				Label:    fmt.Sprintf("%s%s", grade, previouslySelectedMarker),
				//Emoji: &disgord.Emoji{
				//	Name: "📨",
				//},
			},
		)
	}

	err := s.SendInteractionResponse(ctx, h, interactionResponse)

	return err
}

func RespondCommandFailure(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
	message string,
) error {
	err := s.SendInteractionResponse(ctx, h, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Flags: disgord.MessageFlagEphemeral,
			Content: fmt.Sprintf(
				"💥 **BOOM !**\n"+
					"\n"+
					"%s\n"+
					"",
				message,
			),
		},
	})

	return err
}
