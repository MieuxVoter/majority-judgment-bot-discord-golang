package domain

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	db "main/src/database"
	"strings"
)

func RespondWithPollUi(
	input Input,
	poll *db.Poll,
	proposals []*db.Proposal,
	replaceMessage bool,
) error {
	if d, ok := input.(DiscordInput); ok {

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

		err := d.Session.SendInteractionResponse(d.Context, d.Interaction, &disgord.CreateInteractionResponse{
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
	return fmt.Errorf("not supported atm")
}

func RespondWithJudgmentUi(
	input Input,
	judgeSnowflake string,
	proposal *db.Proposal,
	poll *db.Poll,
	previousJudgment *db.Judgment,
	replaceMessage bool,
) error {
	if d, isDiscord := input.(DiscordInput); isDiscord {

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

		return d.Session.SendInteractionResponse(d.Context, d.Interaction, interactionResponse)
	}

	return fmt.Errorf("unsupported vendor")
}

func escapeCsvValue(value string) string {
	return strings.ReplaceAll(value, "\"", "")
}

func RespondBallotsInspection(
	input Input,
	poll *db.Poll,
	proposals []db.Proposal,
	judgments []db.Judgment,
) error {

	csvString := "judge_id"
	for _, proposal := range proposals {
		csvString += fmt.Sprintf(", \"%s\"", escapeCsvValue(proposal.Name))
	}

	currentJudge := ""
	for k, judgment := range judgments {
		if judgment.JudgeSnowflake == currentJudge {
			continue
		}
		currentJudge = judgment.JudgeSnowflake
		csvString += "\n" + currentJudge

		// Extra complexity is to handle missing judgments
		pk := 0
		for _, proposal := range proposals {
			judgmentOfJudge := judgments[k+pk]
			val := "0"
			if judgmentOfJudge.ProposalId == proposal.Id {
				val = fmt.Sprint(judgmentOfJudge.Grade)
				pk += 1
			}
			csvString += ", " + escapeCsvValue(val)
		}
	}

	if d, isDiscord := input.(DiscordInput); isDiscord {

		messageType := disgord.InteractionCallbackChannelMessageWithSource
		csvFile := disgord.CreateMessageFile{
			Reader:   strings.NewReader(csvString),
			FileName: fmt.Sprintf("poll_%d.csv", poll.Id),
		}
		content := fmt.Sprintf("🏛 Here are the individual ballots for the poll **%s** :", poll.Subject)
		interactionResponse := &disgord.CreateInteractionResponse{
			Type: messageType,
			Data: &disgord.CreateInteractionResponseData{
				Flags: disgord.MessageFlagEphemeral,
				Files: []disgord.CreateMessageFile{
					csvFile,
				},
				Content: content,
			},
		}

		return d.Session.SendInteractionResponse(d.Context, d.Interaction, interactionResponse)
	}

	return fmt.Errorf("unsupported vendor")
}

// deprecated
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

func RespondServerError(
	input Input,
	message string,
) error {
	if d, isDiscord := input.(DiscordInput); isDiscord {
		messageType := disgord.InteractionCallbackChannelMessageWithSource
		err := d.Session.SendInteractionResponse(d.Context, d.Interaction, &disgord.CreateInteractionResponse{
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

	return fmt.Errorf("unsupported vendor")
}

func RespondUserError(
	input Input,
	message string,
) error {
	if d, isDiscord := input.(DiscordInput); isDiscord {
		messageType := disgord.InteractionCallbackChannelMessageWithSource
		err := d.Session.SendInteractionResponse(d.Context, d.Interaction, &disgord.CreateInteractionResponse{
			Type: messageType,
			Data: &disgord.CreateInteractionResponseData{
				Flags: disgord.MessageFlagEphemeral,
				Content: fmt.Sprintf(
					"🍄 **Ooops**\n"+
						"\n"+
						"%s\n"+
						"",
					message,
				),
			},
		})

		return err
	}

	return fmt.Errorf("unsupported vendor")
}
