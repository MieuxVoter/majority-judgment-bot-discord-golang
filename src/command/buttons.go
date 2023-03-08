package command

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	db "main/src/database"
	"regexp"
	"strconv"
)

var pollParticipationRegex = regexp.MustCompile("^button_participate:(?P<pollId>\\d+)$")
var pollJudgmentRegex = regexp.MustCompile("^button_judge:(?P<proposalId>\\d+):(?P<gradeLevel>\\d+)$")

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

func HandleButtonParticipate(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
) (handled bool, err error) {
	handled = false
	err = nil

	matches := findNamedMatches(pollParticipationRegex, h.Data.CustomID)
	pollIdAsString, isMatchFound := matches["pollId"]

	if !isMatchFound {
		return
	}

	handled = true

	// Get the judge that clicked the button to participate
	judge := h.Member

	// todo: check the judge's permissions to vote, somehow

	// Get the poll this button is for
	pollId, err := strconv.ParseUint(pollIdAsString, 10, 64)
	if err != nil {
		return false, err
	}
	poll := db.Poll{Id: pollId}
	found, err := db.Orm.Get(&poll)
	if !found {
		err = RespondCommandFailure(ctx, s, h, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get past judgments of the judge on this poll
	judgments, err := db.GetJudgmentsByJudgeOnPoll(db.Orm, judge, &poll)

	// Get the proposals of the poll
	proposals, err := db.GetPollProposals(db.Orm, &poll)
	if err != nil {
		return false, nil
	}

	if len(proposals) == 0 {
		err = RespondCommandFailure(ctx, s, h, "Wait a minute…  This poll has no proposals !?")
		return
	}

	// Shuffle proposals perhaps? todo
	// Pick one proposal (the first)
	proposal := proposals[0]

	var judgment *db.Judgment = nil
	for _, j := range judgments {
		if j.ProposalId == proposal.Id {
			judgment = &j
		}
	}

	// Show the UI to judge that proposal
	err = RespondWithJudgmentUi(ctx, s, h, judge, &proposal, &poll, judgment, false)

	return
}

func HandleButtonJudge(
	ctx context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
) (handled bool, err error) {
	handled = false
	err = nil

	matches := findNamedMatches(pollJudgmentRegex, h.Data.CustomID)
	proposalIdAsString, isProposalIdFound := matches["proposalId"]
	gradeLevelAsString, isGradeLevelFound := matches["gradeLevel"]

	if !isGradeLevelFound || !isProposalIdFound {
		return
	}

	handled = true

	// Get the judge that clicked the button in order to judge
	judge := h.Member

	// todo: check the judge's permissions to judge, somehow

	// Get the proposal this button is for
	proposalId, err := strconv.ParseUint(proposalIdAsString, 10, 64)
	if err != nil {
		return false, err
	}
	proposal := db.Proposal{Id: proposalId}
	found, err := db.Orm.Get(&proposal)
	if !found {
		err = RespondCommandFailure(ctx, s, h, "Oh noes!  This proposal was probably deleted.")
		return
	}
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Ooops.  This proposal was probably deleted.")
		return
	}

	// Get the grade level
	gradeLevel, err := strconv.ParseUint(gradeLevelAsString, 10, 8)
	if err != nil {
		return false, err
	}

	// Get the poll this proposal is attached to
	poll := db.Poll{Id: proposal.PollId}
	found, err = db.Orm.Get(&poll)
	if !found {
		err = RespondCommandFailure(ctx, s, h, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get all past judgments of the judge on this poll
	judgments, err := db.GetJudgmentsByJudgeOnPoll(db.Orm, judge, &poll)
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Nein!.")
		return
	}

	// Get past judgment of this judge on this proposal
	var pastJudgment *db.Judgment = nil
	for _, judgment := range judgments {
		if judgment.ProposalId == proposalId {
			pastJudgment = &judgment
			break
		}
	}

	// Record the judgment, by either updating or inserting
	if pastJudgment != nil {
		oldGrade := pastJudgment.Grade
		pastJudgment.Grade = uint8(gradeLevel)
		updated, err := db.Orm.Update(pastJudgment, &db.Judgment{
			JudgeSnowflake: pastJudgment.JudgeSnowflake,
			ProposalId:     pastJudgment.ProposalId,
			PollId:         pastJudgment.PollId,
			Grade:          oldGrade,
		})
		if updated == 0 {
			return false, fmt.Errorf("did not find a judgment to update")
		}
		if err != nil {
			return false, err
		}
	} else {
		pastJudgment = &db.Judgment{
			JudgeSnowflake: judge.UserID.String(),
			ProposalId:     proposalId,
			PollId:         poll.Id,
			Grade:          uint8(gradeLevel),
		}
		_, err = db.Orm.InsertOne(pastJudgment)
		if err != nil {
			return false, err
		}
	}

	// Get all the proposals of the poll
	proposals, err := db.GetPollProposals(db.Orm, &poll)
	if err != nil {
		return false, err
	}

	amountOfProposals := len(proposals)
	if amountOfProposals == 0 {
		err = RespondCommandFailure(ctx, s, h, "Wait a minute…  This proposal has no proposals !?")
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

	// Shuffle them and pick one todo
	//proposal := proposals[0]

	if nextProposal != nil {

		// Get the past judgment (if any) of this judge on the next proposal
		var nextJudgment *db.Judgment = nil
		for _, j := range judgments {
			if j.ProposalId == nextProposal.Id {
				nextJudgment = &j
			}
		}

		// Show the UI to judge the next proposal
		err = RespondWithJudgmentUi(ctx, s, h, judge, nextProposal, &poll, nextJudgment, true)

	} else {

		message := "Here's the summary of your judgments:\n" +
			"- TODO\n" +
			"- FIXME"
		err = s.SendInteractionResponse(ctx, h, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackUpdateMessage,
			Data: &disgord.CreateInteractionResponseData{
				Flags: disgord.MessageFlagEphemeral | disgord.MessageFlagSupressEmbeds,
				Content: fmt.Sprintf(
					"✅ **CONGRATULATIONS!**"+
						" "+
						"%s\n"+
						"",
					message,
				),
				Embeds: []*disgord.Embed{},
			},
		})
		return

	}

	return
}
