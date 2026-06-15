package domain

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/sarulabs/di/v2"
	"log"
	"log/slog"
	"main/src/container"
	db "main/src/database"
	"main/src/locales"
	"main/src/provider"
	"main/src/security"
	"main/src/services"
	"regexp"
	"strconv"
	"xorm.io/xorm"
)

var buttonPollResultsRegex = regexp.MustCompile("^/button/poll/(?P<pollId>\\d+)/results$")
var buttonPollResultsPattern = "/button/poll/{pollId}/results"

// PollResultsButton is the button the user presses to (privately) see the results
type PollResultsButton struct {
	orm          *xorm.Engine
	gradings     *services.Gradings
	localization *locales.Localization
	logger       *slog.Logger
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

	handled, err = handlePollResult(
		b.orm,
		b.localization,
		b.logger,
		input,
		pollId,
		true,
	)

	return
}

func handlePollResult(
	orm *xorm.Engine,
	localization *locales.Localization,
	logger *slog.Logger,
	input provider.Input,
	pollId uint64,
	asPrivateMessage bool,
) (handled bool, err error) {

	localizer := localization.GetLocalizer(input.GetActorLanguage())
	handled = true

	poll := &db.Poll{Id: pollId}
	hasFoundPoll, err := orm.Get(poll)
	if !hasFoundPoll {
		err = RespondServerError(input, localizer.T("ErrorPollNotFound"))
		return
	}
	if err != nil {
		err = RespondServerError(input, localizer.T("ErrorPollNotFound"))
		return
	}

	proposals, err := db.GetPollProposals(orm, poll)
	if err != nil {
		return false, err
	}
	if len(proposals) == 0 {
		err = RespondServerError(input, localizer.T("ErrorPollHasNoProposals"))
		return
	}

	amountOfJudges, err := db.CountBallots(orm, poll)
	if err != nil {
		err = RespondServerError(input, localizer.T("ErrorPollCannotCountBallots"))
		return
	}

	// Collect the tallies of the proposals
	proposalsTallies := make([]*judgment.ProposalTally, 0, len(proposals))
	for _, proposal := range proposals {
		proposalGradesTally := make([]uint64, 0)
		for gradeLevel := range poll.GetGradingSlice(services.GetGradings()) {
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
		Proposals:      proposalsTallies,
		AmountOfJudges: amountOfJudges,
	}
	//pollTally.GuessAmountOfJudges() // nope ; we count them with an SQL query instead

	// Rule: the "worst" grade is the default grade
	err = pollTally.BalanceWithStaticDefault(0)
	if err != nil {
		return
	}

	// Rule: proposals are ranked in the merit profile
	deliberator := &judgment.MajorityJudgment{}
	pollResult, err := deliberator.Deliberate(pollTally)
	if err != nil {
		return
	}

	// Rule: GTFO if there are no judgments
	if pollTally.AmountOfJudges == 0 {
		message := localizer.T("ErrorNoParticipants")
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

	content := fmt.Sprintf("🏆 ")
	if len(winnersSlice) > 1 {
		content += fmt.Sprintf(
			`_%s_ %s`,
			localizer.T("TheFavoriteProposalsAre"),
			winners,
		)
	} else {
		content += fmt.Sprintf(
			`_%s_ %s`,
			localizer.T("TheFavoriteProposalIs"),
			winners,
		)
	}

	judgeVendorId, _ := input.GetActorVendorId()
	canInspect, _ := security.CanUserInspectBallots(orm, judgeVendorId, poll)

	logger.Info(
		"showing poll results",
		"id", pollId,
		"title", title,
		"winners", winners,
	)

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
				orm:          ctn.Get("database.engine").(*xorm.Engine),
				gradings:     ctn.Get("gradings").(*services.Gradings),
				localization: ctn.Get("localization").(*locales.Localization),
				logger:       ctn.Get("logger").(*slog.Logger),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("button.poll.result failed to build", err)
	}
}
