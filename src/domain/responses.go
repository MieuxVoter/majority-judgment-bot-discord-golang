package domain

import (
	"fmt"
	"github.com/andersfylling/disgord"
	db "main/src/database"
	"main/src/provider"
	"main/src/security"
	"strings"
)

func RespondPollView(
	input provider.Input,
	poll *db.Poll,
	proposals []*db.Proposal,
	replaceMessage bool,
) error {
	return provider.GetResponder(input).RespondPollView(
		input,
		poll,
		proposals,
		replaceMessage,
	)
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
