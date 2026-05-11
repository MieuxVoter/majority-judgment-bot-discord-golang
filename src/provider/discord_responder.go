package provider

import (
	"bytes"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/mieuxvoter/merit-profile-library-go/merit"
	"github.com/mskrha/svg2png"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/security"
	"time"
	"xorm.io/xorm"
)

//func getInteractionEvent(input Input) (*handler.InteractionEvent, error) {
//	discordCommandInput, isDiscordCommandInput := (input).(DiscordCommandInput)
//	if isDiscordCommandInput {
//		interactionEvent, isInteractionEvent := (*discordCommandInput.Event).(*handler.InteractionEvent)
//		if isInteractionEvent {
//			return &interactionEvent, nil
//		}
//	}
//	//_, isDiscord = (input).(DiscordButtonInput)
//}

// DiscordResponder implements provider.ResponderInterface for Discord
type DiscordResponder struct {
	orm *xorm.Engine
}

func (r DiscordResponder) sanitizeTitle(title string) string {
	return security.TruncateString(title, 256)
}

//func (r DiscordResponder) convertButtonField(field *provider.ButtonField) *disgord.MessageComponent {
//	component := &disgord.MessageComponent{
//		Type:     disgord.MessageComponentButton,
//		Style:    field.Style,
//		Label:    field.Label,
//		CustomID: field.Id,
//	}
//
//	if field.Url != "" {
//		component.Url = field.Url
//		component.Style = disgord.Link
//		component.CustomID = ""
//	}
//
//	if field.Emote != "" {
//		component.Emoji = &disgord.Emoji{
//			Name: field.Emote,
//		}
//	}
//
//	return component
//}

func (r DiscordResponder) Matches(input Input) bool {
	_, isDiscord := (input).(DiscordCommandInput)
	if !isDiscord {
		_, isDiscord = (input).(DiscordButtonInput)
	}
	return isDiscord
}

func (r DiscordResponder) RespondWithMessage(input Input, message string, ephemeral bool) error {
	if d, isDiscord := (input).(DiscordInteraction); isDiscord {
		msg := discord.MessageCreate{
			Content: message,
		}

		if ephemeral {
			msg = msg.WithFlags(discord.MessageFlagEphemeral)
		}

		return d.CreateMessage(msg)
	}

	return RaiseInvalidProviderError("Discord:RespondWithMessage")
}

//func (r DiscordResponder) RespondWithMessageAndImage(
//	input provider.Input,
//	message string,
//	imageUrl string,
//	ephemeral bool,
//) error {
//	if d, isDiscord := (input).(provider.DiscordCommandInput); isDiscord {
//		response := &disgord.CreateInteractionResponse{
//			Type: disgord.InteractionCallbackChannelMessageWithSource,
//			Data: &disgord.CreateInteractionResponseData{
//				Content: message,
//				Embeds: []*disgord.Embed{
//					{
//						Type: disgord.EmbedTypeImage,
//						//Title: title,
//						Image: &disgord.EmbedImage{
//							URL: imageUrl,
//						},
//					},
//				},
//			},
//		}
//		if ephemeral {
//			response.Data.Flags |= disgord.MessageFlagEphemeral
//		}
//
//		return d.Session.SendInteractionResponse(d.Context, d.Interaction, response)
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondWithMessage")
//}

//func (r DiscordResponder) RespondWithMessageAndButtons(
//	input provider.Input,
//	message string,
//	buttons []*provider.ButtonField,
//	ephemeral bool,
//) error {
//	if d, isDiscord := (input).(provider.DiscordCommandInput); isDiscord {
//		response := &disgord.CreateInteractionResponse{
//			Type: disgord.InteractionCallbackChannelMessageWithSource,
//			Data: &disgord.CreateInteractionResponseData{
//				Content: message,
//			},
//		}
//
//		if ephemeral {
//			response.Data.Flags |= disgord.MessageFlagEphemeral
//		}
//
//		if len(buttons) > 0 {
//			row := &disgord.MessageComponent{
//				Type:       disgord.MessageComponentActionRow,
//				CustomID:   "message_action_row",
//				Components: make([]*disgord.MessageComponent, 0),
//			}
//
//			for _, button := range buttons {
//				row.Components = append(row.Components, r.convertButtonField(button))
//			}
//
//			response.Data.Components = make([]*disgord.MessageComponent, 0)
//			response.Data.Components = append(response.Data.Components, row)
//		}
//
//		return d.Session.SendInteractionResponse(d.Context, d.Interaction, response)
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondWithMessage")
//}

func (r DiscordResponder) RespondPollView(
	input Input,
	poll *db.Poll,
	proposals []*db.Proposal,
	replaceMessage bool,
) error {
	if d, isDiscord := input.(DiscordCommandInput); isDiscord {

		title := "### " + r.sanitizeTitle(poll.Subject)
		amountOfVotes, _ := db.CountBallots(r.orm, poll)

		description := ""
		if len(proposals) > 0 {
			for _, proposal := range proposals {
				description += "- "
				proposalName := security.RemoveMarkdown(proposal.Name)
				description += security.TruncateEllipsis(proposalName, 256)
				description += "\n"
			}
		}

		msg := discord.MessageCreate{
			Flags: discord.MessageFlagIsComponentsV2,
			Components: []discord.LayoutComponent{
				discord.NewContainer(
					discord.NewSection(
						discord.NewTextDisplay(title),
						discord.NewTextDisplay(description),
					).WithAccessory(
						discord.NewPrimaryButton(
							"Vote",
							fmt.Sprintf("/button/poll/%d/vote", poll.Id),
						).WithEmoji(
							discord.ComponentEmoji{Name: "🗳"},
						),
					),
					discord.NewSmallSeparator(),
					discord.NewSection(
						discord.NewTextDisplay(
							fmt.Sprintf("%d votes", amountOfVotes),
						),
					).WithAccessory(
						discord.NewSecondaryButton(
							"Show results",
							fmt.Sprintf("/button/poll/%d/results", poll.Id),
						).WithEmoji(
							discord.ComponentEmoji{Name: "🔎"},
							//discord.ComponentEmoji{Name: "🏆"},
						),
					),
				),
			},
		}

		return d.Event.CreateMessage(msg)
	}

	return RaiseInvalidProviderError("Discord:RespondPollView")
}

func (r DiscordResponder) RespondWithJudgmentUi(
	input Input,
	proposal *db.Proposal,
	poll *db.Poll,
	previousJudgment *db.Judgment,
	replaceMessage bool,
) error {
	if d, isDiscord := input.(DiscordInteraction); isDiscord {

		title := "### " + security.TruncateString(proposal.Name, 256)

		gradeButtons := make([]discord.InteractiveComponent, 0, 5)
		for gradeLevel, grade := range poll.GetGradingSlice() {
			customId := fmt.Sprintf("/button/poll/%d/judge/%d/as/%d", poll.Id, proposal.Id, gradeLevel)
			emoji := discord.ComponentEmoji{Name: grade}
			button := discord.NewSecondaryButton("", customId).WithEmoji(emoji)
			if previousJudgment != nil && previousJudgment.Grade == uint8(gradeLevel) {
				button = discord.NewSuccessButton("", customId).WithEmoji(emoji)
			}
			gradeButtons = append(gradeButtons, button)
		}

		flags := discord.MessageFlagIsComponentsV2 | discord.MessageFlagEphemeral
		components := []discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay(title),
				discord.NewActionRow(gradeButtons...),
			),
		}

		if replaceMessage {
			return d.UpdateMessage(discord.MessageUpdate{
				// We need pointers here, not sure why
				Flags:      &flags,
				Components: &components,
			})
		}

		return d.CreateMessage(discord.MessageCreate{
			Flags:      flags,
			Components: components,
		})
	}

	return RaiseInvalidProviderError("Discord:RespondWithJudgmentUi")
}

func (r DiscordResponder) RespondBallotSummary(
	input Input,
	poll *db.Poll,
	proposals []db.Proposal,
	judgments []db.Judgment,
) error {
	if d, isDiscord := input.(DiscordInteraction); isDiscord {

		title := "### ✅ **A VOTÉ**"
		message := "Here's the summary of your judgments:"
		summary := ""
		for k := range judgments {
			icon := poll.GetGradeIcon(judgments[k].Grade)
			summary += fmt.Sprintf("- %s ⋅ %s\n", icon, proposals[k].Name)
		}

		flags := discord.MessageFlagIsComponentsV2 | discord.MessageFlagEphemeral
		components := []discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay(title),
				discord.NewTextDisplay(message),
				discord.NewTextDisplay(summary),
			),
		}

		return d.UpdateMessage(discord.MessageUpdate{
			Flags:      &flags,
			Components: &components,
		})
	}

	return RaiseInvalidProviderError("Discord:RespondBallotSummary")
}

func (r DiscordResponder) RespondPollResult(
	input Input,
	poll *db.Poll,
	proposals []db.Proposal,
	pollTally *judgment.PollTally,
	pollResult *judgment.PollResult,
	title string,
	message string,
	asPrivateMessage bool,
	canInspect bool,
) error {
	if d, isDiscord := input.(DiscordInteraction); isDiscord {

		rendererProposals := make([]merit.Proposal, len(proposals))
		for i, proposal := range pollResult.ProposalsSorted {
			rendererProposals[i] = merit.Proposal{
				Name:  proposals[proposal.Index].Name,
				Tally: pollTally.Proposals[proposal.Index].Tally,
			}
		}
		svg, err := merit.RenderLinearProfileSVG(
			rendererProposals,
			// fake data for testing
			//[]merit.Proposal{
			//	{
			//		Name:  "Jonlukz",
			//		Tally: []uint64{1, 2, 2, 4, 3},
			//	},
			//	{
			//		Name:  "Bourdella",
			//		Tally: []uint64{5, 3, 3, 0, 1},
			//	},
			//	{
			//		Name:  "Rempaillot",
			//		Tally: []uint64{3, 4, 4, 1, 0},
			//	},
			//},
		)

		if err != nil {
			return err
		}

		//fmt.Println(svg)
		svgBytes := []byte(svg)

		// Discord does not render SVG files (although it's somewhat safe in img tags)
		// So we need to create a raster version of our merit profile.
		// To that effect, we use svg2png which in turn uses inkscape internally.
		// It's not pretty, but it works.  Docker will help a little.
		converter := svg2png.New()

		// We can also tell it where our inkscape binary resides.
		// I've compiled inkscape myself; you probably won't need this.
		//err = converter.SetBinary("/usr/local/bin/inkscape")
		//if err != nil {
		//	fmt.Println(err)
		//	return err
		//}

		pngBytes, err := converter.Convert(svgBytes)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Debug dump of the generated PNG file.
		//err = os.WriteFile("merit.png", pngBytes, 0644)
		//if err != nil {
		//	fmt.Println(err)
		//	return err
		//}

		// TODO: slugify the poll subject, truncate it and append it
		imageFilename := fmt.Sprintf(
			`merit-profile-%s`,
			time.Now().Format("20060102"),
		)
		imageDescription := fmt.Sprintf(
			`Merit Profile of the poll: %s`,
			poll.Subject,
		)

		participantsPluralization := ""
		if pollTally.AmountOfJudges > 1 {
			participantsPluralization = "s"
		}

		flags := discord.MessageFlagIsComponentsV2
		if asPrivateMessage {
			flags |= discord.MessageFlagEphemeral
		}

		msg := discord.MessageCreate{
			Flags: flags,
			Components: []discord.LayoutComponent{
				discord.NewContainer(
					discord.NewTextDisplay("### "+title),
					discord.NewTextDisplay(message),
					discord.NewSmallSeparator(),
					discord.NewMediaGallery(
						discord.MediaGalleryItem{
							Media: discord.UnfurledMediaItem{
								URL: "attachment://" + imageFilename + ".png",
							},
							Description: "A merit profile",
							Spoiler:     false,
						},
					),
					// Shows a link to download the file. (not what we want)
					//discord.NewFileComponent(
					//	"attachment://"+imageFilename+".png",
					//),
					discord.NewSmallSeparator(),
					discord.NewSection(
						discord.NewTextDisplay(
							fmt.Sprintf(
								`%d participant%s`,
								pollTally.AmountOfJudges,
								participantsPluralization,
							),
						),
					).WithAccessory(
						discord.NewSecondaryButton(
							"Publish",
							fmt.Sprintf("/button/poll/%d/publish", poll.Id),
						).WithEmoji(
							discord.ComponentEmoji{Name: "📢"},
						),
					),
				),
			},
			Files: []*discord.File{
				discord.NewFile(
					imageFilename+".png",
					imageDescription,
					bytes.NewReader(pngBytes),
					discord.FileFlagsNone,
				),
				//discord.NewFile(
				//	imageFilename+".svg",
				//	imageDescription,
				//	bytes.NewReader(svgBytes),
				//	discord.FileFlagsNone,
				//),
			},
		}
		return d.CreateMessage(msg)

		//if canInspect {
		//	response.Data.Components[0].Components = append(
		//		response.Data.Components[0].Components,
		//		&disgord.MessageComponent{
		//			Type:  disgord.MessageComponentButton,
		//			Style: disgord.Secondary,
		//			Label: "Inspect Ballots",
		//			Emoji: &disgord.Emoji{
		//				Name: "🕵",
		//			},
		//			CustomID: fmt.Sprintf("button_inspect:%d", poll.Id),
		//		},
		//	)
		//}
	}

	return RaiseInvalidProviderError("Discord:RespondPollResult")
}

//func (r DiscordResponder) RespondBallotsInspection(
//	input provider.Input,
//	poll *db.Poll,
//	proposals []db.Proposal,
//	judgments []db.Judgment,
//) error {
//
//	csvString := "judge_id"
//	for _, proposal := range proposals {
//		csvString += fmt.Sprintf(", \"%s\"", security.EscapeCsvValue(proposal.Name))
//	}
//
//	currentJudgeVendorId := ""
//	for k, jt := range judgments {
//		if jt.JudgeSnowflake == currentJudgeVendorId {
//			continue
//		}
//		currentJudgeVendorId = jt.JudgeSnowflake
//		csvString += fmt.Sprintf("\n\"%s\"", currentJudgeVendorId)
//
//		// Extra complexity is to handle missing judgments
//		pk := 0
//		for _, proposal := range proposals {
//			judgmentOfJudge := judgments[k+pk]
//			val := "0"
//			if judgmentOfJudge.ProposalId == proposal.Id {
//				val = fmt.Sprint(judgmentOfJudge.Grade)
//				pk += 1
//			}
//			csvString += ", " + security.EscapeCsvValue(val)
//		}
//	}
//
//	if d, isDiscord := input.(provider.DiscordCommandInput); isDiscord {
//
//		messageType := disgord.InteractionCallbackChannelMessageWithSource
//		csvFile := disgord.CreateMessageFile{
//			Reader:   strings.NewReader(csvString),
//			FileName: fmt.Sprintf("poll_%d.csv", poll.Id),
//		}
//		content := fmt.Sprintf("🏛 Here are the individual ballots for the poll **%s** :", poll.Subject)
//		interactionResponse := &disgord.CreateInteractionResponse{
//			Type: messageType,
//			Data: &disgord.CreateInteractionResponseData{
//				Flags: disgord.MessageFlagEphemeral,
//				Files: []disgord.CreateMessageFile{
//					csvFile,
//				},
//				Content: content,
//			},
//		}
//
//		return d.Session.SendInteractionResponse(d.Context, d.Interaction, interactionResponse)
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondBallotInspection")
//}

func (r DiscordResponder) RespondServerError(
	input Input,
	message string,
) error {
	if _, isDiscord := input.(DiscordInteraction); isDiscord {
		return r.RespondWithMessage(
			input,
			fmt.Sprintf(
				"### 💥 **BOOM !**\n"+
					"\n"+
					"%s\n",
				message,
			),
			true,
		)
	}

	return RaiseInvalidProviderError("Discord:RespondServerError")
}

func (r DiscordResponder) RespondUserError(
	input Input,
	message string,
) error {
	if _, isDiscord := input.(DiscordInteraction); isDiscord {
		return r.RespondWithMessage(
			input,
			fmt.Sprintf(
				"### 🍄 **Ooopsie !**\n"+
					"\n"+
					"%s\n"+
					"",
				message,
			),
			true,
		)
	}

	return RaiseInvalidProviderError("Discord:RespondUserError")
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "responder.discord",
		Build: func(ctn di.Container) (interface{}, error) {
			responder := &DiscordResponder{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return responder, nil
		},
	})

	if err != nil {
		log.Fatalln("service responder.discord failed to build", err)
	}
}
