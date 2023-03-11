package command

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	db "main/src/database"
	"main/src/security"
	"main/src/services"
	"net/url"
	"regexp"
	"strconv"
)

var pollParticipationRegex = regexp.MustCompile("^button_participate:(?P<pollId>\\d+)$")
var pollDeliberationRegex = regexp.MustCompile("^button_deliberate:(?P<pollId>\\d+)$")
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
	found, err := db.GetEngine().Get(&poll)
	if !found {
		err = RespondCommandFailure(ctx, s, h, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get past judgments of the judge on this poll
	judgments, err := db.GetJudgmentsByJudgeOnPoll(db.GetEngine(), judge, &poll)

	// Get the proposals of the poll
	proposals, err := db.GetPollProposals(db.GetEngine(), &poll)
	if err != nil {
		return false, nil
	}

	if len(proposals) == 0 {
		err = RespondCommandFailure(ctx, s, h, "Wait a minute…  This poll has no proposals !?")
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
	err = RespondWithJudgmentUi(ctx, s, h, judge, &proposal, &poll, pastJudgment, false)

	return
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

	// Get the guild this poll is for
	guild, err := db.GetGuild(db.GetEngine(), h.GuildID.String())
	if err != nil {
		services.GetLogger().Errorln(err)
		err = RespondCommandFailure(ctx, s, h, "Oh snap!  This guild is not registered.")
		return
	}

	// Check if the guild is banned or not
	canParticipate, err := security.CanGuildParticipate(db.GetEngine(), guild)
	if !canParticipate {
		err = RespondCommandFailure(ctx, s, h, "Ouch!  This guild is banned.")
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
	found, err = db.GetEngine().Get(&poll)
	if !found {
		err = RespondCommandFailure(ctx, s, h, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get all past judgments of the judge on this poll
	judgments, err := db.GetJudgmentsByJudgeOnPoll(db.GetEngine(), judge, &poll)
	if err != nil {
		err = RespondCommandFailure(ctx, s, h, "Nein!")
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
		// /!. This does not update when gradeLevel is zero, unless Cols() is specified
		// > When this param is the pointer of struct, only non-empty and non-zero field will be updated to database.
		// > from https://xorm.io/docs/chapter-06/readme/
		updated, err := db.GetEngine().Cols("grade").Update(pastJudgment, &db.Judgment{
			JudgeSnowflake: pastJudgment.JudgeSnowflake,
			ProposalId:     pastJudgment.ProposalId,
			PollId:         pastJudgment.PollId,
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
		_, err = db.GetEngine().InsertOne(pastJudgment)
		if err != nil {
			return false, err
		}
	}

	// Get all the proposals of the poll
	proposals, err := db.GetPollProposals(db.GetEngine(), &poll)
	if err != nil {
		return false, err
	}

	amountOfProposals := len(proposals)
	if amountOfProposals == 0 {
		err = RespondCommandFailure(ctx, s, h, "Wait a minute…  This poll has no proposals !?  Go :fish:")
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
		err = RespondWithJudgmentUi(ctx, s, h, judge, nextProposal, &poll, nextJudgment, true)

	} else {

		summary := ""
		for k := range judgments {
			if k > 0 {
				summary += "  —  "
			}

			icon := poll.GetGradeIcon(judgments[k].Grade)
			summary += fmt.Sprintf("(%s %s %s)", icon, proposals[k].Name, icon)
		}
		message := "Here's the summary of your judgments:\n" + summary
		err = s.SendInteractionResponse(ctx, h, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackUpdateMessage,
			Data: &disgord.CreateInteractionResponseData{
				//Flags: disgord.MessageFlagEphemeral | disgord.MessageFlagSupressEmbeds,
				Flags: disgord.MessageFlagEphemeral,
				//Content: fmt.Sprintf(
				//	"✅ **WELL DONE!**"+
				//		" "+
				//		"%s\n"+
				//		"",
				//	message,
				//),
				Embeds: []*disgord.Embed{
					{
						Title:       "✅ **WELL DONE!**",
						Description: message,
					},
				},
			},
		})
	}

	return
}
