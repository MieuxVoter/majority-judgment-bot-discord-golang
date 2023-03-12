package command

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
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

	handled = true

	// todo: perhaps check the actor's permissions to deliberate?

	// Get the poll this button is for
	pollId, err := strconv.ParseUint(pollIdAsString, 10, 64)
	if err != nil {
		return false, err
	}
	poll := &db.Poll{Id: pollId}
	foundPoll, err := service.orm.Get(poll)
	if !foundPoll {
		err = RespondServerError(input, "Oh noes!  This poll was probably deleted.")
		return
	}
	if err != nil {
		err = RespondServerError(input, "Ooops.  This poll was probably deleted.")
		return
	}

	// Get the proposals of the poll
	proposals, err := db.GetPollProposals(db.GetEngine(), poll)
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
		// fixme: strip markdown from proposal names?
		winners += fmt.Sprintf("**%s**", proposal.Name)
	}

	if d, isDiscord := input.(DiscordInput); isDiscord {

		content := fmt.Sprintf(
			"🤖 _Here are the most plebiscited proposals:_ %s",
			winners,
		)
		err = d.Session.SendInteractionResponse(d.Context, d.Interaction, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Flags:   disgord.MessageFlagEphemeral,
				Content: content,

				Embeds: []*disgord.Embed{
					{
						//Title:       "",
						Type: disgord.EmbedTypeImage,
						//Description: "",
						//URL:         "",
						//Timestamp:   disgord.Time{},
						//Color:       0,
						//Footer:      nil,
						Image: &disgord.EmbedImage{
							URL: imageUrl,
							//ProxyURL: "",
							//Height:   0,
							//Width:    0,
						},
						//Thumbnail:   nil,
						//Video:       nil,
						//Provider:    nil,
						//Author:      nil,
						//Fields:      nil,
					},
				},

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
	} else {
		log.Fatalln("vendor not supported")
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
