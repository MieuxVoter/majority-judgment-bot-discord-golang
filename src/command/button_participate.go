package command

import (
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"regexp"
	"strconv"
	"xorm.io/xorm"
)

type ParticipateButton struct {
	orm *xorm.Engine
}

func (b ParticipateButton) Handle(input Input) (bool, error) {
	return handleButtonParticipate(&b, input)
}

var pollParticipationRegex = regexp.MustCompile("^button_participate:(?P<pollId>\\d+)$")

func handleButtonParticipate(
	button *ParticipateButton,
	input Input,
) (handled bool, err error) {
	handled = false
	err = nil

	buttonName, err := input.GetButtonName()
	if err != nil {
		return
	}
	matches := findNamedMatches(pollParticipationRegex, buttonName)
	pollIdAsString, isMatchFound := matches["pollId"]

	if !isMatchFound {
		return
	}

	handled = true

	// Get the judge that clicked the button to participate
	judgeSnowflake, err := input.GetActorVendorId()
	if err != nil {
		err = RespondServerError(input, "Ooops.  _I can't figure out who you are._")
		return
	}

	// todo: check the judge's permissions to vote, somehow

	// Get the poll this button is for
	pollId, err := strconv.ParseUint(pollIdAsString, 10, 64)
	if err != nil {
		return false, err
	}
	poll := db.Poll{Id: pollId}
	found, err := button.orm.Get(&poll)
	if !found {
		err = RespondUserError(input, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondUserError(input, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get past judgments of the judge on this poll
	judgments, err := db.GetJudgmentsByJudgeOnPoll(button.orm, judgeSnowflake, &poll)

	// Get the proposals of the poll
	proposals, err := db.GetPollProposals(button.orm, &poll)
	if err != nil {
		return false, nil
	}

	if len(proposals) == 0 {
		err = RespondUserError(input, "Wait a minute…  This poll has no proposals !?")
		return
	}

	// Shuffle proposals perhaps?

	// Pick one proposal (the first)
	proposal := proposals[0]

	// Collect the past judgment (if any) on this proposal by this judge
	var pastJudgment *db.Judgment = nil
	for _, j := range judgments {
		if j.ProposalId == proposal.Id {
			pastJudgment = &j
			break
		}
	}

	// Show the UI to judge that proposal
	err = RespondWithJudgmentUi(input, judgeSnowflake, &proposal, &poll, pastJudgment, false)

	return
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "button.participate",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &ParticipateButton{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("button.participate failed to build", err)
	}
}
