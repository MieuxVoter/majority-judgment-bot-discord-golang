package domain

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"xorm.io/xorm"
)

type InfoCommand struct {
	orm *xorm.Engine
}

func (c InfoCommand) Define() *disgord.ApplicationCommandOption {
	return &disgord.ApplicationCommandOption{
		Name:        "info",
		Description: "Display miscellaneous information about this bot on this server",
		Type:        disgord.OptionTypeSubCommand,
	}
}

func (c InfoCommand) Matches(command string) bool {
	return command == "info"
}

func (c InfoCommand) Handle(input Input) (handled bool, err error) {
	if d, ok := (input).(DiscordInput); ok {
		return true, handleInfoCommand(c, d)
	}
	return false, fmt.Errorf("unknown vendor")
}

func handleInfoCommand(
	command InfoCommand,
	input DiscordInput,
) error {
	guildVendorId, _ := input.GetGuildVendorId()
	guild, err := db.GetOrCreateGuild(command.orm, guildVendorId)
	if err != nil {
		message := "Could not access the guild.  _Suddenly, everything is on fire. 🔥_"
		return RespondUserError(input, message)
	}

	err = input.Session.SendInteractionResponse(input.Context, input.Interaction, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Flags: disgord.MessageFlagEphemeral,
			Content: "" +
				"🤖🗩 _Here is some information about myself on this server._\n" +
				"\n" +
				"Polls remaining: " + fmt.Sprintf("%d", guild.Quota) +
				"\n" +
				"",
		},
	})

	return err
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command.info",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &InfoCommand{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("command.info failed to build", err)
	}
}
