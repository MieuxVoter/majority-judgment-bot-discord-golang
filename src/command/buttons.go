package command

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	db "main/src/database"
	"net/url"
	"regexp"
	"strconv"
)

var pollDeliberationRegex = regexp.MustCompile("^button_deliberate:(?P<pollId>\\d+)$")

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

type Button interface {
	Handle(input Input) (bool, error)
}

func HandleButtonDeliberate(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
) (handled bool, err error) {
	handled = false
	err = nil

	matches := findNamedMatches(pollDeliberationRegex, h.Data.CustomID)
	pollIdAsString, isMatchFound := matches["pollId"]

	if !isMatchFound {
		return
	}

	handled = true

	// Get the judge that clicked the button to participate
	//actor := h.Member

	// todo: check the actor's permissions to deliberate, somehow

	// Get the poll this button is for
	pollId, err := strconv.ParseUint(pollIdAsString, 10, 64)
	if err != nil {
		return false, err
	}
	poll := &db.Poll{Id: pollId}
	foundPoll, err := db.GetEngine().Get(poll)
	if !foundPoll {
		err = RespondCommandFailure(ctx, s, h, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get the proposals of the poll
	proposals, err := db.GetPollProposals(db.GetEngine(), poll)
	if err != nil {
		return false, err
	}
	if len(proposals) == 0 {
		err = RespondCommandFailure(ctx, s, h, "Wait a minute…  This poll has no proposals !?")
		return
	}

	// Rank the proposals
	proposalsTallies := make([]*judgment.ProposalTally, 0, len(proposals))
	for _, proposal := range proposals {
		proposalGradesTally := make([]uint64, 0)
		for gradeLevel := range poll.GetGradingSlice() {
			gradeAmount, errCount := db.CountGrades(db.GetEngine(), poll, &proposal, uint8(gradeLevel))
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
	for proposalResultIndex, proposalResult := range result.ProposalsSorted {
		if proposalResult.Rank > 1 {
			break
		}
		proposal := proposals[proposalResult.Index]
		if proposalResultIndex > 0 {
			winners += fmt.Sprintf(", ")
		}
		winners += fmt.Sprintf("%s", proposal.Name)
	}

	content := fmt.Sprintf(
		"🤖 _Here are the results:_ **%s**\n"+
			"%s",
		winners, imageUrl,
	)
	err = s.SendInteractionResponse(ctx, h, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Flags:   disgord.MessageFlagEphemeral,
			Content: content,

			//Attachments: []*disgord.Attachment{
			//	{
			//		ID:         0,
			//		Filename:   imageUrl,
			//		//Size:       0,
			//		URL:        imageUrl,
			//		//ProxyURL:   "",
			//		//Height:     0,
			//		//Width:      0,
			//		//SpoilerTag: true,
			//	},
			//},
		},
	})
	if err != nil {
		return
	}

	//imageUrl := "https://oas.mieuxvoter.fr/%s.png?subject=%s&proposals[]=HAHA&proposals[]=HIHI"

	//// Get all judgments emitted on this poll
	//judgments, err := db.GetJudgmentsOnPoll(db.GetEngine(), poll)

	// Shuffle proposals perhaps? todo
	// Pick one proposal (the first)
	//proposal := proposals[0]

	//var judgment *db.Judgment = nil
	//for _, j := range judgments {
	//	if j.ProposalId == proposal.Id {
	//		judgment = &j
	//		break
	//	}
	//}

	// Show the UI to judge that proposal
	//err = RespondWithJudgmentUi(ctx, s, h, judge, &proposal, &poll, judgment, false)

	return
}
