package domain

import (
	"fmt"
	"github.com/andersfylling/disgord"
	db "main/src/database"
	"main/src/provider"
	"main/src/security"
	"strings"
)

func RespondWithPollUi(
	input provider.Input,
	poll *db.Poll,
	proposals []*db.Proposal,
	replaceMessage bool,
) error {
	if d, ok := input.(provider.DiscordInput); ok {

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
	input provider.Input,
	proposal *db.Proposal,
	poll *db.Poll,
	previousJudgment *db.Judgment,
	replaceMessage bool,
) error {
	return provider.GetResponder(input).RespondWithJudgmentUi(
		input,
		proposal,
		poll,
		previousJudgment,
		replaceMessage,
	)
}

func RespondBallotsInspection(
	input provider.Input,
	poll *db.Poll,
	proposals []db.Proposal,
	judgments []db.Judgment,
) error {

	csvString := "judge_id"
	for _, proposal := range proposals {
		csvString += fmt.Sprintf(", \"%s\"", security.EscapeCsvValue(proposal.Name))
	}

	currentJudge := ""
	for k, judgment := range judgments {
		if judgment.JudgeSnowflake == currentJudge {
			continue
		}
		currentJudge = judgment.JudgeSnowflake
		csvString += fmt.Sprintf("\n\"%s\"", currentJudge)

		// Extra complexity is to handle missing judgments
		pk := 0
		for _, proposal := range proposals {
			judgmentOfJudge := judgments[k+pk]
			val := "0"
			if judgmentOfJudge.ProposalId == proposal.Id {
				val = fmt.Sprint(judgmentOfJudge.Grade)
				pk += 1
			}
			csvString += ", " + security.EscapeCsvValue(val)
		}
	}

	if d, isDiscord := input.(provider.DiscordInput); isDiscord {

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

func RespondServerError(input provider.Input, message string) error {
	return provider.GetResponder(input).RespondServerError(input, message)
}

func RespondUserError(input provider.Input, message string) error {
	return provider.GetResponder(input).RespondUserError(input, message)
}
