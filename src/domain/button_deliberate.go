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

var pollDeliberationRegex = regexp.MustCompile("^button_deliberate:(?P<pollId>\\d+)$")

type DeliberateButton struct {
	orm *xorm.Engine
}

func (service DeliberateButton) Handle(input provider.Input) (handled bool, err error) {

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
	input provider.Input,
	pollId uint64,
	asPrivateMessage bool,
) (handled bool, err error) {

	handled = true

	// todo: perhaps check the actor's permissions to deliberate? (use a middleware!)

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
	pollResult, errDelib := deliberator.Deliberate(pollTally)
	if nil != errDelib {
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
			winners += fmt.Sprintf(", ")
		}
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
			"_Here are the most consensual proposals:_ %s",
			winners,
		)
	} else {
		content += fmt.Sprintf(
			"_Here is the most consensual proposal:_ %s",
			winners,
		)
	}

	judgeVendorId, _ := input.GetActorVendorId()
	canInspect, _ := security.CanUserInspectBallots(orm, judgeVendorId, poll)

	err = provider.GetResponder(input).RespondDeliberation(
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
