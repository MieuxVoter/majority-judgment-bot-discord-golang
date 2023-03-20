package discord

import (
	"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	"main/src/provider"
)

// Responder implements provider.ResponderInterface
type Responder struct{}

func (r Responder) Matches(input provider.Input) bool {
	_, isDiscord := (input).(provider.DiscordInput)
	return isDiscord
}

func (r Responder) RespondWithMessage(input provider.Input, message string, ephemeral bool) error {
	if d, isDiscord := (input).(provider.DiscordInput); isDiscord {
		response := &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: message,
			},
		}
		if ephemeral {
			response.Data.Flags |= disgord.MessageFlagEphemeral
		}

		return d.Session.SendInteractionResponse(d.Context, d.Interaction, response)
	}

	return provider.RaiseInvalidProviderError("Discord:RespondWithMessage")
}

func init() {
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
