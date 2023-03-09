package command

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	db "main/src/database"
)

func RespondWithPollUi(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
	//judge *disgord.Member,
	poll *db.Poll,
	proposals []*db.Proposal,
	//previousJudgment *db.Judgment,
	replaceMessage bool,
) error {
	messageType := disgord.InteractionCallbackChannelMessageWithSource
	if replaceMessage {
		messageType = disgord.InteractionCallbackUpdateMessage
	}

	pollEmbedHero := &disgord.Embed{
		Title: fmt.Sprintf("⚖ `#%d` %s", poll.Id, poll.Subject),
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

	err := s.SendInteractionResponse(ctx, h, &disgord.CreateInteractionResponse{
		Type: messageType,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				pollEmbedHero,
			},
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
						{
							Type:     disgord.MessageComponentButton,
							Style:    disgord.Secondary,
							CustomID: fmt.Sprintf("button_deliberate:%d", poll.Id),
							Label:    "View Results",
							Emoji: &disgord.Emoji{
								Name: "🔎",
							},
						},
					},
				},
			},
		},
	})

	return err
}

func RespondWithJudgmentUi(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
	judge *disgord.Member,
	proposal *db.Proposal,
	poll *db.Poll,
	previousJudgment *db.Judgment,
	replaceMessage bool,
) error {
	messageType := disgord.InteractionCallbackChannelMessageWithSource
	if replaceMessage {
		messageType = disgord.InteractionCallbackUpdateMessage
	}
	interactionResponse := &disgord.CreateInteractionResponse{
		Type: messageType,
		Data: &disgord.CreateInteractionResponseData{
			Flags: disgord.MessageFlagEphemeral,
			Embeds: []*disgord.Embed{
				{
					Title:       fmt.Sprintf("⚖ `#%d` %s", poll.Id, proposal.Name),
					Description: fmt.Sprintf("What do you think of **_%s_** as _%s_ ?", proposal.Name, poll.Subject),
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
	messageType := disgord.InteractionCallbackChannelMessageWithSource
	err := s.SendInteractionResponse(ctx, h, &disgord.CreateInteractionResponse{
		Type: messageType,
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
