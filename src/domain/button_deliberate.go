package domain

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/security"
	"net/url"
	"regexp"
	"strconv"
	"xorm.io/xorm"
)

var pollDeliberationRegex = regexp.MustCompile("^button_deliberate:(?P<pollId>\\d+)$")

type DeliberateButton struct {
	orm *xorm.Engine
}

func (service DeliberateButton) Handle(input Input) (handled bool, err error) {

	handled = false
	err = nil

	buttonName, err := input.GetButtonName()
	if err != nil {
		return
	}

	matches := findNamedMatches(pollDeliberationRegex, buttonName)
	pollIdAsString, isMatchFound := matches["pollId"]

	if !isMatchFound {
		return
	}

	pollId, errParse := strconv.ParseUint(pollIdAsString, 10, 64)
	if errParse != nil {
		return false, errParse
	}

	handled, err = handleDeliberation(service.orm, input, pollId, true)

	return
}

func handleDeliberation(
	orm *xorm.Engine,
	input Input,
	pollId uint64,
	asPrivateMessage bool,
) (handled bool, err error) {

	handled = true

	// todo: perhaps check the actor's permissions to deliberate?

	// Get the poll this button is for
	poll := &db.Poll{Id: pollId}
	hasFoundPoll, err := orm.Get(poll)
	if !hasFoundPoll {
		err = RespondServerError(input, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondServerError(input, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get the proposals of the poll
	proposals, err := db.GetPollProposals(orm, poll)
	if err != nil {
		return false, err
	}
	if len(proposals) == 0 {
		err = RespondServerError(input, "Wait a minute…  This poll has no proposals !?")
		return
	}

	// Rank the proposals
	proposalsTallies := make([]*judgment.ProposalTally, 0, len(proposals))
	for _, proposal := range proposals {
		proposalGradesTally := make([]uint64, 0)
		for gradeLevel := range poll.GetGradingSlice() {
			gradeAmount, errCount := db.CountGrades(orm, poll, &proposal, uint8(gradeLevel))
			if errCount != nil {
				return false, errCount
			}

			proposalGradesTally = append(proposalGradesTally, gradeAmount)
		}
		proposalTally := &judgment.ProposalTally{Tally: proposalGradesTally}
		proposalsTallies = append(proposalsTallies, proposalTally)
	}

	pollTally := &judgment.PollTally{
		Proposals: proposalsTallies,
	}
	pollTally.GuessAmountOfJudges()
	err = pollTally.BalanceWithStaticDefault(0)
	if err != nil {
		return
	}
	deliberator := &judgment.MajorityJudgment{}
	result, err := deliberator.Deliberate(pollTally)

	if nil != err {
		return
	}

	// GTFO if there are no judgments
	if pollTally.AmountOfJudges == 0 {
		message := "There are no participants to this poll.  Please try again when the poll has had participants."
		err = RespondUserError(input, message)
	}

	// Build the path to the merit profile image
	fileNameNoExt := ""
	query := fmt.Sprintf("?subject=%s", url.QueryEscape(poll.Subject))
	for proposalResultIndex, proposalResult := range result.ProposalsSorted {
		proposal := proposals[proposalResult.Index]

		if proposalResultIndex > 0 {
			fileNameNoExt += "_"
		}
		for gradeLevel := range poll.GetGradingSlice() {
			gradeAmount := pollTally.Proposals[proposalResult.Index].Tally[gradeLevel]

			if gradeLevel > 0 {
				fileNameNoExt += "-"
			}
			fileNameNoExt += fmt.Sprintf("%d", gradeAmount)
		}
		query += fmt.Sprintf("&proposals[]=%s", url.QueryEscape(proposal.Name))
	}

	oasDomain := "https://oas.mieuxvoter.fr"
	imageUrl := fmt.Sprintf(
		"%s/%s.png%s",
		oasDomain, fileNameNoExt, query,
	)

	winners := ""
	winnersSlice := make([]string, 0)
	for proposalResultIndex, proposalResult := range result.ProposalsSorted {
		if proposalResult.Rank > 1 {
			break
		}
		proposal := proposals[proposalResult.Index]
		if proposalResultIndex > 0 {
			winners += fmt.Sprintf(", ")
		}
		// fixme: strip markdown from proposal names?
		winners += fmt.Sprintf("**%s**", proposal.Name)
		winnersSlice = append(winnersSlice, proposal.Name)
	}

	title := fmt.Sprintf("%d participant", pollTally.AmountOfJudges)
	if pollTally.AmountOfJudges > 1 {
		title += "s"
	}

	content := fmt.Sprintf("🤖⚖  ")
	if len(winnersSlice) > 1 {
		content += fmt.Sprintf(
			"_Here are the most plebiscited proposals:_ %s",
			winners,
		)
	} else {
		content += fmt.Sprintf(
			"_Here is the most plebiscited proposal:_ %s",
			winners,
		)
	}

	if d, isDiscord := input.(DiscordInput); isDiscord {

		judgeVendorId, _ := input.GetActorVendorId()
		canInspect, _ := security.CanUserInspectBallots(orm, judgeVendorId, poll)

		response := &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: content,
				Flags:   0,
				Embeds: []*disgord.Embed{
					{
						Type:  disgord.EmbedTypeImage,
						Title: title,
						Image: &disgord.EmbedImage{
							URL: imageUrl,
						},
					},
				},
			},
		}
		if asPrivateMessage || canInspect {
			response.Data.Components = []*disgord.MessageComponent{
				{
					Type:       disgord.MessageComponentActionRow,
					CustomID:   "deliberate_action_row",
					Components: []*disgord.MessageComponent{},
				},
			}
		}
		if asPrivateMessage {
			//response.Type = disgord.InteractionCallbackChannelMessageWithSource
			response.Data.Flags |= disgord.MessageFlagEphemeral

			response.Data.Components[0].Components = append(
				response.Data.Components[0].Components,
				&disgord.MessageComponent{
					Type:  disgord.MessageComponentButton,
					Style: disgord.Primary,
					Label: "Publish",
					Emoji: &disgord.Emoji{
						Name: "📢",
					},
					CustomID: fmt.Sprintf("button_publish:%d", poll.Id),
				},
			)
		} else {
			response.Data.Flags |= disgord.MessageFlagSourceMessageDeleted
		}

		if canInspect {
			response.Data.Components[0].Components = append(
				response.Data.Components[0].Components,
				&disgord.MessageComponent{
					Type:  disgord.MessageComponentButton,
					Style: disgord.Secondary,
					Label: "Inspect Ballots",
					Emoji: &disgord.Emoji{
						Name: "🕵",
					},
					CustomID: fmt.Sprintf("button_inspect:%d", poll.Id),
				},
			)
		}

		err = d.Session.SendInteractionResponse(d.Context, d.Interaction, response)
		if err != nil {
			return
		}
	} else {
		log.Fatalln("vendor not supported")
	}

	return
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "button.deliberate",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &DeliberateButton{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("button.deliberate failed to build", err)
	}
}
