package domain

import (
	//"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/provider"
	"xorm.io/xorm"
)

const FeedbackCommandSlug = "feedback"

type FeedbackCommand struct {
	orm *xorm.Engine
}

func (c FeedbackCommand) GetEmote() string {
	return "💡"
}

func (c FeedbackCommand) GetName() string {
	return FeedbackCommandSlug
}

func (c FeedbackCommand) GetDescription() string {
	return "Send a message to my creators"
}

//func (c FeedbackCommand) DefineForDiscord() *disgord.ApplicationCommandOption {
//	return &disgord.ApplicationCommandOption{
//		Name:        c.GetName(),
//		Description: c.GetEmote() + " " + c.GetDescription(),
//		Type:        disgord.OptionTypeSubCommand,
//		Options: []*disgord.ApplicationCommandOption{
//			{
//				Type:        disgord.OptionTypeString,
//				Required:    true,
//				Name:        "message",
//				Description: "The message you wish to send us — intolerance won't be tolerated",
//			},
//		},
//	}
//}

func (c FeedbackCommand) Matches(command string) bool {
	return command == FeedbackCommandSlug
}

func (c FeedbackCommand) Handle(input provider.Input) (handled bool, err error) {
	return true, handleFeedbackCommand(c.orm, input)
}

func handleFeedbackCommand(orm *xorm.Engine, input provider.Input) error {

	messageSent, _ := input.GetOption(FeedbackCommandSlug, "message", "")
	if messageSent == "" {
		message := "Please provide a message with your feedback."
		return RespondUserError(input, message)
	}

	guildVendorId, _ := input.GetGuildVendorId()
	guild, err := db.GetGuild(orm, guildVendorId)
	if err != nil {
		message := "This community never held any poll.  Feedback requires experience."
		return RespondUserError(input, message)
	}

	actorVendorId, err := input.GetActorVendorId()
	if err != nil {
		message := "I cannot figure out who you are.  Feedback was canceled."
		return RespondUserError(input, message)
	}
	actorName, _ := input.GetActorName()

	feedback := &db.Feedback{
		GuildId:        guild.Id,
		Content:        messageSent,
		AuthorVendorId: actorVendorId,
		AuthorName:     actorName,
	}
	_, err = orm.InsertOne(feedback)
	if err != nil {
		message := "I cannot write into my memory anymore.  Please contact my creator in any other way."
		return RespondServerError(input, message)
	}

	message := "🤖🗩 _Your feedback was successfully recorded.  **Thank you !**_" +
		"\n" +
		""

	return provider.GetResponder(input).RespondWithMessage(input, message, true)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command." + FeedbackCommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &FeedbackCommand{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})

	if err != nil {
		log.Fatalf("command.%s failed to build : %s\n", FeedbackCommandSlug, err)
	}
}
