package domain

import (
	"fmt"
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/provider"
	"main/src/security"
	"main/src/services"
	"regexp"
	"strconv"
	"xorm.io/xorm"
)

var buttonPollJudgeRegex = regexp.MustCompile("^/button/poll/(?P<pollId>\\d+)/judge/(?P<proposalId>\\d+)/as/(?P<gradeLevel>\\d+)$")
var buttonPollJudgePattern = "/button/poll/{pollId}/judge/{proposalId}/as/{gradeLevel}"

// PollJudgeButton is the button users press when grading a single proposal
type PollJudgeButton struct {
	orm *xorm.Engine
}

func (b PollJudgeButton) GetRegex() *regexp.Regexp {
	return buttonPollJudgeRegex
}

func (b PollJudgeButton) GetPattern() string {
	return buttonPollJudgePattern
}

func (b PollJudgeButton) Handle(
	input provider.ButtonInput,
) (handled bool, err error) {

	handled = false
	err = nil

	buttonName, err := input.GetButtonName()
	if err != nil {
		return
	}
	matches := findNamedMatches(buttonPollJudgeRegex, buttonName)
	proposalIdAsString, isProposalIdFound := matches["proposalId"]
	gradeLevelAsString, isGradeLevelFound := matches["gradeLevel"]

	if !isGradeLevelFound || !isProposalIdFound {
		return
	}

	handled = true

	// Get the judge that clicked the button in order to judge the proposal
	judgeVendorId, err := input.GetActorVendorId()
	if err != nil {
		err = RespondServerError(input, "Ooops.  _I can't figure out who you are._")
		return
	}

	// Get the guild this poll is for
	guildVendorId, err := input.GetGuildVendorId()
	guild, err := db.GetGuild(b.orm, guildVendorId)
	if err != nil {
		services.GetLogger().Errorln(err)
		err = RespondServerError(input, "Oh snap!  This guild is not registered.")
		return
	}

	// Check if the guild is banned or not
	canParticipate, err := security.CanGuildParticipate(b.orm, guild)
	if !canParticipate {
		err = RespondUserError(input, "Ouch!  This guild is banned.")
		return
	}

	// Get the proposal this button is for
	proposalId, err := strconv.ParseUint(proposalIdAsString, 10, 64)
	if err != nil {
		return false, err
	}
	proposal := db.Proposal{Id: proposalId}
	found, err := db.GetEngine().Get(&proposal)
	if !found {
		err = RespondUserError(input, "Oh noes!  This proposal was probably deleted.")
		return
	}
	if err != nil {
		err = RespondUserError(input, "Ooops.  This proposal was probably deleted.")
		return
	}

	// Get the grade level
	gradeLevel, err := strconv.ParseUint(gradeLevelAsString, 10, 8)
	if err != nil {
		return false, err
	}

	// Get the poll this proposal is attached to
	poll := &db.Poll{Id: proposal.PollId}
	found, err = b.orm.Get(poll)
	if !found {
		err = RespondUserError(input, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondUserError(input, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get all past judgments of the judge on this poll
	judgments, err := db.GetJudgmentsByJudgeOnPoll(b.orm, judgeVendorId, poll)
	if err != nil {
		err = RespondServerError(input, "_Nein!_ "+err.Error())
		return
	}

	// Get past judgment of this judge on this proposal
	var pastJudgment *db.Judgment = nil
	for k := range judgments {
		if judgments[k].ProposalId == proposalId {
			pastJudgment = &(judgments[k])
			break
		}
	}

	// Record the judgment, by either updating or inserting
	if pastJudgment != nil {
		pastJudgment.Grade = uint8(gradeLevel)
		// /!. This does not update when gradeLevel is zero, unless Cols() is specified (thankfully, it is).
		// > When this param is the pointer of struct, only non-empty and non-zero field will be updated to database.
		// > from https://xorm.io/docs/chapter-06/readme/
		updated, errUpdate := b.orm.Cols("grade").Update(pastJudgment, &db.Judgment{
			JudgeSnowflake: pastJudgment.JudgeSnowflake,
			ProposalId:     pastJudgment.ProposalId,
		})
		if updated == 0 {
			return false, fmt.Errorf("did not find a judgment to update")
		}
		if errUpdate != nil {
			return false, errUpdate
		}

	} else {
		newJudgment := &db.Judgment{
			JudgeSnowflake: judgeVendorId,
			ProposalId:     proposalId,
			Grade:          uint8(gradeLevel),
		}
		_, err = b.orm.InsertOne(newJudgment)
		if err != nil {
			return false, err
		}
		judgments = append(judgments, *newJudgment)
	}

	// Get all the proposals of the poll
	proposals, err := db.GetPollProposals(b.orm, poll)
	if err != nil {
		return false, err
	}

	amountOfProposals := len(proposals)
	if amountOfProposals == 0 {
		err = RespondServerError(input, "Wait a minute…  This poll has no proposals !?  Go :fish:")
		return
	}

	// Get the next proposal to go to after this one (if any)
	var nextProposal *db.Proposal = nil
	useNextProposal := false
	for _, p := range proposals {
		if useNextProposal {
			nextProposal = &p
			break
		}
		if p.Id == proposal.Id {
			useNextProposal = true
		}
	}

	// Now we either show the next proposal, or the conclusion (a summary of judgments)
	if nextProposal != nil {

		// Get the past judgment (if any) of this judge on the next proposal
		var nextJudgment *db.Judgment = nil
		for _, j := range judgments {
			if j.ProposalId == nextProposal.Id {
				nextJudgment = &j
				break
			}
		}

		// Show the UI to judge the next proposal
		err = RespondWithJudgmentUi(input, nextProposal, poll, nextJudgment, true)
		return

	} else {
		// Rule: When all proposals are judged, show a summary of emitted judgments.
		err = RespondBallotSummary(
			input,
			poll,
			proposals,
			judgments,
		)
		return
	}
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "button.poll.judge",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &PollJudgeButton{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("button.poll.judge failed to build", err)
	}
}
