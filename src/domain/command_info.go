package domain

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/provider"
	"main/src/security"
	"xorm.io/xorm"
)

// InfoCommand displays miscellaneous information about the bot
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

func (c InfoCommand) Handle(input provider.Input) (handled bool, err error) {
	return true, handleInfoCommand(c, input)
}

func handleInfoCommand(
	command InfoCommand,
	input provider.Input,
) error {
	guildVendorId, _ := input.GetGuildVendorId()
	guild, err := db.GetOrCreateGuild(command.orm, guildVendorId)
	if err != nil {
		message := "Could not access the guild.  _Suddenly, everything is on fire._ 🔥"
		return RespondUserError(input, message)
	}
	allPollsAmount, errCountAll := db.CountPolls(command.orm)
	if errCountAll != nil {
		message := "Could not count the polls.  _Suddenly, Notre-Dame is on fire._ 🔥"
		return RespondUserError(input, message)
	}
	guildPollsAmount, errCountGuildPolls := db.CountGuildPolls(command.orm, guild)
	if errCountGuildPolls != nil {
		message := "Could not count this guild's polls.  _Suddenly, Australia is on fire._ 🔥"
		return RespondUserError(input, message)
	}

	message := "" +
		"🤖🗩 _Here is some information about myself._\n" +
		"\n" +
		"Total amount of polls by this community" + fmt.Sprintf(" : `%d`\n", guildPollsAmount) +
		"Total amount of polls across all communities" + fmt.Sprintf(" : `%d`\n", allPollsAmount) +
		"Remaining polls' quota of this community" + fmt.Sprintf(" : `%d`\n", guild.Quota) +
		"Version" + fmt.Sprintf(" : `%s`\n", security.GetVersion()) +
		"Guild Identifier" + fmt.Sprintf(" : `%s`\n", guild.Snowflake) +
		"\n" +
		""

	return provider.GetResponder(input).RespondWithMessage(input, message, true)
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
