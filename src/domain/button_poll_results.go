package domain

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/provider"
	"main/src/security"
	"regexp"
	"strconv"
	"xorm.io/xorm"
)

var buttonPollResultsRegex = regexp.MustCompile("^/button/poll/(?P<pollId>\\d+)/results$")
var buttonPollResultsPattern = "/button/poll/{pollId}/results"

type PollResultsButton struct {
	orm *xorm.Engine
}

func (b PollResultsButton) GetRegex() *regexp.Regexp {
	return buttonPollResultsRegex
}

func (b PollResultsButton) GetPattern() string {
	return buttonPollResultsPattern
}

func (b PollResultsButton) Handle(
	input provider.ButtonInput,
) (handled bool, err error) {

	handled = false
	err = nil

	buttonName, err := input.GetButtonName()
	if err != nil {
		return
	}

	matches := findNamedMatches(b.GetRegex(), buttonName)
	pollIdAsString, isMatchFound := matches["pollId"]

	if !isMatchFound {
		return
	}

	pollId, errParse := strconv.ParseUint(pollIdAsString, 10, 64)
	if errParse != nil {
		return false, errParse
	}

	handled, err = handlePollResult(b.orm, input, pollId, true)

	return
}

func handlePollResult(
	orm *xorm.Engine,
	input provider.Input,
	pollId uint64,
	asPrivateMessage bool,
) (handled bool, err error) {

	handled = true

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

	proposals, err := db.GetPollProposals(orm, poll)
	if err != nil {
		return false, err
	}
	if len(proposals) == 0 {
		err = RespondServerError(input, "Wait a minute…  This poll has no proposals !?")
		return
	}

	// Rule: Proposals are ranked in the merit profile
	proposalsTallies := make([]*judgment.ProposalTally, 0, len(proposals))
	for _, proposal := range proposals {
		proposalGradesTally := make([]uint64, 0)
		for gradeLevel := range poll.GetGradingSlice() {
			gradeAmount, errCount := db.CountGradesReceived(orm, poll, &proposal, uint8(gradeLevel))
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
	pollTally.GuessAmountOfJudges() // TBD: we should count them with an SQL query
	err = pollTally.BalanceWithStaticDefault(0)
	if err != nil {
		return
	}

	deliberator := &judgment.MajorityJudgment{}
	pollResult, err := deliberator.Deliberate(pollTally)
	if err != nil {
		return
	}

	// Rule: GTFO if there are no judgments
	if pollTally.AmountOfJudges == 0 {
		message := "There are no participants to this poll.  Please try again when the poll has had participants."
		err = RespondUserError(input, message)
		return
	}

	winners := ""
	winnersSlice := make([]string, 0)
	for proposalResultIndex, proposalResult := range pollResult.ProposalsSorted {
		if proposalResult.Rank > 1 {
			break
		}
		proposal := proposals[proposalResult.Index]
		if proposalResultIndex > 0 {
			winners += fmt.Sprintf(" `💠` ")
		}
		proposalName := security.RemoveMarkdown(proposal.Name)
		proposalName = security.TruncateEllipsis(proposalName, 256)
		winners += fmt.Sprintf("**%s**", proposalName)
		winnersSlice = append(winnersSlice, proposal.Name)
	}

	title := poll.Subject
	//title := fmt.Sprintf("%d participant", pollTally.AmountOfJudges)
	//if pollTally.AmountOfJudges > 1 {
	//	title += "s"
	//}

	content := fmt.Sprintf("⚖  ")
	if len(winnersSlice) > 1 {
		content += fmt.Sprintf(
			`_The most consensual proposals are_ %s`,
			winners,
		)
	} else {
		content += fmt.Sprintf(
			`_The most consensual proposal is_ %s`,
			winners,
		)
	}

	judgeVendorId, _ := input.GetActorVendorId()
	canInspect, _ := security.CanUserInspectBallots(orm, judgeVendorId, poll)

	err = provider.GetResponder(input).RespondPollResult(
		input,
		poll,
		proposals,
		pollTally,
		pollResult,
		title,
		content,
		asPrivateMessage,
		canInspect,
	)

	return
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "button.poll.result",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &PollResultsButton{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("button.poll.result failed to build", err)
	}
}
