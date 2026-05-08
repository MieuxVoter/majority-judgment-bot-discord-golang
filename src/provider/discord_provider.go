package provider

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	"main/src/security"
)

// Responder implements provider.ResponderInterface for Discord
type Responder struct{}

func (r Responder) sanitizeTitle(title string) string {
	return security.TruncateString(title, 256)
}

//func (r Responder) convertButtonField(field *provider.ButtonField) *disgord.MessageComponent {
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

func (r Responder) Matches(input Input) bool {
	_, isDiscord := (input).(DiscordInput)
	return isDiscord
}

func (r Responder) RespondWithMessage(input Input, message string, ephemeral bool) error {
	if d, isDiscord := (input).(DiscordInput); isDiscord {
		msg := discord.MessageCreate{
			Content: message,
		}

		//if ephemeral {
		//	msg = msg.WithFlags(discord.MessageFlagEphemeral)
		//}

		return d.Event.CreateMessage(msg)
	}

	return RaiseInvalidProviderError("Discord:RespondWithMessage")
}

//func (r Responder) RespondWithMessageAndImage(
//	input provider.Input,
//	message string,
//	imageUrl string,
//	ephemeral bool,
//) error {
//	if d, isDiscord := (input).(provider.DiscordInput); isDiscord {
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

//func (r Responder) RespondWithMessageAndButtons(
//	input provider.Input,
//	message string,
//	buttons []*provider.ButtonField,
//	ephemeral bool,
//) error {
//	if d, isDiscord := (input).(provider.DiscordInput); isDiscord {
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

//func (r Responder) RespondPollView(
//	input provider.Input,
//	poll *db.Poll,
//	proposals []*db.Proposal,
//	replaceMessage bool,
//) error {
//	if d, isDiscord := input.(provider.DiscordInput); isDiscord {
//
//		title := fmt.Sprintf("⚖ `#%d` %s", poll.Id, poll.Subject)
//		title = r.sanitizeTitle(title)
//
//		messageType := disgord.InteractionCallbackChannelMessageWithSource
//		if replaceMessage {
//			messageType = disgord.InteractionCallbackUpdateMessage
//		}
//
//		pollEmbedHero := &disgord.Embed{
//			Title: title,
//		}
//		if len(proposals) > 0 {
//			description := ""
//			for i, proposal := range proposals {
//				if i > 0 {
//					description += " **|** "
//				}
//
//				proposalName := security.RemoveMarkdown(proposal.Name)
//				description += security.TruncateEllipsis(proposalName, 256)
//			}
//			pollEmbedHero.Description = description
//		}
//
//		response := &disgord.CreateInteractionResponse{
//			Type: messageType,
//			Data: &disgord.CreateInteractionResponseData{
//				Embeds: []*disgord.Embed{
//					pollEmbedHero,
//				},
//				Components: []*disgord.MessageComponent{
//					{
//						Type:     disgord.MessageComponentActionRow,
//						CustomID: "poll_action_row",
//						Components: []*disgord.MessageComponent{
//							{
//								Type:     disgord.MessageComponentButton,
//								Style:    disgord.Success,
//								CustomID: fmt.Sprintf("button_participate:%d", poll.Id),
//								Label:    "Participate",
//								Emoji: &disgord.Emoji{
//									Name: "📨",
//								},
//							},
//							{
//								Type:     disgord.MessageComponentButton,
//								Style:    disgord.Secondary,
//								CustomID: fmt.Sprintf("button_deliberate:%d", poll.Id),
//								Label:    "View Results",
//								Emoji: &disgord.Emoji{
//									Name: "🔎",
//								},
//							},
//						},
//					},
//				},
//			},
//		}
//
//		return d.Session.SendInteractionResponse(d.Context, d.Interaction, response)
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondPollView")
//}

//func (r Responder) RespondWithJudgmentUi(
//	input provider.Input,
//	proposal *db.Proposal,
//	poll *db.Poll,
//	previousJudgment *db.Judgment,
//	replaceMessage bool,
//) error {
//	if d, isDiscord := input.(provider.DiscordInput); isDiscord {
//
//		title := fmt.Sprintf("⚖ `#%d` %s", poll.Id, proposal.Name)
//		title = security.TruncateString(title, 256)
//		messageType := disgord.InteractionCallbackChannelMessageWithSource
//		if replaceMessage {
//			messageType = disgord.InteractionCallbackUpdateMessage
//		}
//		interactionResponse := &disgord.CreateInteractionResponse{
//			Type: messageType,
//			Data: &disgord.CreateInteractionResponseData{
//				Flags: disgord.MessageFlagEphemeral,
//				Embeds: []*disgord.Embed{
//					{
//						Title:       title,
//						Description: fmt.Sprintf("What do you think of **_%s_** as _%s_ ?", proposal.Name, poll.Subject),
//					},
//				},
//				Components: []*disgord.MessageComponent{
//					{
//						Type:       disgord.MessageComponentActionRow,
//						CustomID:   "poll_action_row",
//						Components: []*disgord.MessageComponent{}, // filled below
//					},
//				},
//			},
//		}
//
//		for gradeLevel, grade := range poll.GetGradingSlice() {
//
//			previouslySelectedMarker := ""
//			if previousJudgment != nil {
//				if uint8(gradeLevel) == previousJudgment.Grade {
//					previouslySelectedMarker = " ✅"
//				}
//			}
//			interactionResponse.Data.Components[0].Components = append(
//				interactionResponse.Data.Components[0].Components,
//				&disgord.MessageComponent{
//					Type:     disgord.MessageComponentButton,
//					Style:    disgord.Primary,
//					CustomID: fmt.Sprintf("button_judge:%d:%d", proposal.Id, gradeLevel),
//					Label:    fmt.Sprintf("%s%s", grade, previouslySelectedMarker),
//				},
//			)
//		}
//
//		return d.Session.SendInteractionResponse(d.Context, d.Interaction, interactionResponse)
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondWithJudgmentUi")
//}

//func (r Responder) RespondJudgmentSummary(
//	input provider.Input,
//	poll *db.Poll,
//	proposals []db.Proposal,
//	judgments []db.Judgment,
//	replaceMessage bool,
//) error {
//	if d, isDiscord := input.(provider.DiscordInput); isDiscord {
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

//func (r Responder) RespondDeliberation(
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
//	if d, isDiscord := input.(provider.DiscordInput); isDiscord {
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

//func (r Responder) RespondBallotsInspection(
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
//	if d, isDiscord := input.(provider.DiscordInput); isDiscord {
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

//func (r Responder) RespondServerError(
//	input provider.Input,
//	message string,
//) error {
//	if d, isDiscord := input.(provider.DiscordInput); isDiscord {
//		messageType := disgord.InteractionCallbackChannelMessageWithSource
//		err := d.Session.SendInteractionResponse(d.Context, d.Interaction, &disgord.CreateInteractionResponse{
//			Type: messageType,
//			Data: &disgord.CreateInteractionResponseData{
//				Flags: disgord.MessageFlagEphemeral,
//				Content: fmt.Sprintf(
//					"💥 **BOOM !**\n"+
//						"\n"+
//						"%s\n"+
//						"",
//					message,
//				),
//			},
//		})
//
//		return err
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondServerError")
//}

//func (r Responder) RespondUserError(
//	input provider.Input,
//	message string,
//) error {
//	if d, isDiscord := input.(provider.DiscordInput); isDiscord {
//		messageType := disgord.InteractionCallbackChannelMessageWithSource
//		err := d.Session.SendInteractionResponse(d.Context, d.Interaction, &disgord.CreateInteractionResponse{
//			Type: messageType,
//			Data: &disgord.CreateInteractionResponseData{
//				Flags: disgord.MessageFlagEphemeral,
//				Content: fmt.Sprintf(
//					"🍄 **Ooops**\n"+
//						"\n"+
//						"%s\n"+
//						"",
//					message,
//				),
//			},
//		})
//
//		return err
//	}
//
//	return provider.RaiseInvalidProviderError("Discord:RespondUserError")
//}

func init() {
	fmt.Println("init() responder.discord")

	err := container.GetBuilder().Add(di.Def{
		Name: "responder.discord",
		Build: func(ctn di.Container) (interface{}, error) {
			responder := &Responder{}
			return responder, nil
		},
	})

	if err != nil {
		log.Fatalln("service responder.discord failed to build", err)
	}
}
