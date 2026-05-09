package provider

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/security"
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
							fmt.Sprintf("button_poll_results:%d", poll.Id),
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

		msg := discord.MessageCreate{
			Flags: discord.MessageFlagIsComponentsV2 | discord.MessageFlagEphemeral,
			Components: []discord.LayoutComponent{
				discord.NewContainer(
					discord.NewTextDisplay(title),
					discord.NewActionRow(gradeButtons...),
				),
			},
		}

		return d.CreateMessage(msg)
	}

	return RaiseInvalidProviderError("Discord:RespondWithJudgmentUi")
}

//func (r DiscordResponder) RespondJudgmentSummary(
//	input provider.Input,
//	poll *db.Poll,
//	proposals []db.Proposal,
//	judgments []db.Judgment,
//	replaceMessage bool,
//) error {
//	if d, isDiscord := input.(provider.DiscordCommandInput); isDiscord {
//
//		summary := ""
//		for k := range judgments {
//			if k > 0 {
//				summary += "  —  "
//			}
//
//			icon := poll.GetGradeIcon(judgments[k].Grade)
//			summary += fmt.Sprintf("(%s %s %s)", icon, proposals[k].Name, icon)
//		}
//		title := "✅ **WELL DONE!**"
//		message := "Here's the summary of your judgments:\n" + summary
//
//		return d.Session.SendInteractionResponse(d.Context, d.Interaction, &disgord.CreateInteractionResponse{
//			Type: disgord.InteractionCallbackUpdateMessage,
//			Data: &disgord.CreateInteractionResponseData{
//				Flags: disgord.MessageFlagEphemeral, // | disgord.MessageFlagSupressEmbeds,
//				Embeds: []*disgord.Embed{
//					{
//						Title:       title,
//						Description: message,
//					},
//				},
//			},
//		})
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondJudgmentSummary")
//}

//func (r DiscordResponder) RespondDeliberation(
//	input provider.Input,
//	poll *db.Poll,
//	proposals []db.Proposal,
//	pollTally *judgment.PollTally,
//	pollResult *judgment.PollResult,
//	title string,
//	message string,
//	asPrivateMessage bool,
//	canInspect bool,
//) error {
//	if d, isDiscord := input.(provider.DiscordCommandInput); isDiscord {
//
//		// Generate the merit profile image URL
//		imageUrl, errImg := network.GetOas().GetMeritProfileUrl(
//			poll,
//			proposals,
//			pollTally,
//			pollResult,
//			"png",
//			MaxUrlLength,
//		)
//		if errImg != nil {
//			imageUrl = ""
//		}
//		imageUrlSvg, errImgSvg := network.GetOas().GetMeritProfileUrl(
//			poll,
//			proposals,
//			pollTally,
//			pollResult,
//			"svg",
//			MaxUrlLength,
//		)
//		if errImgSvg != nil {
//			imageUrlSvg = ""
//		}
//
//		response := &disgord.CreateInteractionResponse{
//			Type: disgord.InteractionCallbackChannelMessageWithSource,
//			Data: &disgord.CreateInteractionResponseData{
//				Content: message,
//				Flags:   0,
//				Embeds: []*disgord.Embed{
//					{
//						Type:  disgord.EmbedTypeImage,
//						Title: title,
//						Image: &disgord.EmbedImage{
//							// Rule: SVG is NOT allowed here, it appears
//							// Rule: 256 characters max
//							URL: imageUrl,
//						},
//					},
//				},
//			},
//		}
//
//		if asPrivateMessage || canInspect {
//			response.Data.Components = []*disgord.MessageComponent{
//				{
//					Type:       disgord.MessageComponentActionRow,
//					CustomID:   "deliberate_action_row",
//					Components: []*disgord.MessageComponent{},
//				},
//			}
//		}
//		if asPrivateMessage {
//			//response.Type = disgord.InteractionCallbackChannelMessageWithSource
//			response.Data.Flags |= disgord.MessageFlagEphemeral
//
//			response.Data.Components[0].Components = append(
//				response.Data.Components[0].Components,
//				&disgord.MessageComponent{
//					Type:  disgord.MessageComponentButton,
//					Style: disgord.Primary,
//					Label: "Publish",
//					Emoji: &disgord.Emoji{
//						Name: "📢",
//					},
//					CustomID: fmt.Sprintf("button_publish:%d", poll.Id),
//				},
//			)
//			if imageUrlSvg != "" {
//				response.Data.Components[0].Components = append(
//					response.Data.Components[0].Components,
//					&disgord.MessageComponent{
//						Type:  disgord.MessageComponentButton,
//						Style: disgord.Link,
//						Label: "As SVG",
//						Emoji: &disgord.Emoji{
//							Name: "✨",
//						},
//						Url: imageUrlSvg,
//					},
//				)
//			}
//		} else {
//			response.Data.Flags |= disgord.MessageFlagSourceMessageDeleted
//		}
//
//		if canInspect {
//			response.Data.Components[0].Components = append(
//				response.Data.Components[0].Components,
//				&disgord.MessageComponent{
//					Type:  disgord.MessageComponentButton,
//					Style: disgord.Secondary,
//					Label: "Inspect Ballots",
//					Emoji: &disgord.Emoji{
//						Name: "🕵",
//					},
//					CustomID: fmt.Sprintf("button_inspect:%d", poll.Id),
//				},
//			)
//		}
//
//		return d.Session.SendInteractionResponse(d.Context, d.Interaction, response)
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondDeliberation")
//}

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
	if _, isDiscord := input.(DiscordCommandInput); isDiscord {
		//return GetResponder(input).RespondWithMessage(
		return r.RespondWithMessage(
			input,
			fmt.Sprintf(
				"💥 **BOOM !**\n"+
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
	if _, isDiscord := input.(DiscordCommandInput); isDiscord {
		return r.RespondWithMessage(
			input,
			fmt.Sprintf(
				"🍄 **Ooops**\n"+
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
