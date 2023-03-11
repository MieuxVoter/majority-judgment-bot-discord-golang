package command

import (
	"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"strconv"
	"strings"
	"xorm.io/xorm"
)

type RerunCommand struct {
	orm *xorm.Engine
}

func (c RerunCommand) Define() *disgord.ApplicationCommandOption {
	return &disgord.ApplicationCommandOption{
		Name:        "rerun",
		Description: "Rerun a copy of a past poll",
		Type:        disgord.OptionTypeSubCommand,
		Options: []*disgord.ApplicationCommandOption{
			{
				Name:        "poll",
				Description: "The poll's numerical identifier, shown after the ⚖",
				Type:        disgord.OptionTypeString,
			},
		},
	}
}

func (c RerunCommand) Matches(command string) bool {
	return command == "rerun"
}

func (c RerunCommand) Handle(input Input) (handled bool, err error) {
	return true, handleRerunCommand(
		c.orm,
		input,
	)
}

func handleRerunCommand(
	orm *xorm.Engine,
	input Input,
) error {

	guildVendorId, _ := input.GetGuildVendorId()
	guild, err := db.GetGuild(orm, guildVendorId)
	if err != nil {
		message := "This guild has no polls to rerun, and is not even registered yet.  " +
			"Please create a poll first, using `/mj create …`."
		return RespondCommandUserError(input, message)
	}

	pollIdString, _ := input.GetOption("rerun", "poll", "")
	if pollIdString == "" {
		// TODO: fetch most recent poll on this guild
		return RespondCommandUserError(input, "Fetching most recent poll is not implemented yet.  "+
			"Please provide a poll identifier.")
	}
	pollIdString = strings.Trim(pollIdString, "#")

	pollId, errConv := strconv.Atoi(pollIdString)
	if errConv != nil {
		message := "🐉 The specified poll identifier is not a number.  " +
			"Please use the _numerical_ identifier of the poll you want to rerun, " +
			"like so `/mj rerun poll:42`."
		return RespondCommandUserError(input, message)
	}
	poll, errPoll := db.FindPoll(orm, uint64(pollId))
	if errPoll != nil {
		return errPoll
	}
	if poll == nil {
		message := "The specified poll was not found.  " +
			"Sorry.  _Here, have a banana instead : 🍌_"
		return RespondCommandUserError(input, message)
	}
	if poll.GuildId != guild.Id {
		message := "The specified poll belongs to another community.  " +
			"No can do!  _Dura lex, sed lex 🏛_"
		return RespondCommandUserError(input, message)
	}

	proposals, errProp := db.GetPollProposals(orm, poll)
	if errProp != nil {
		return errProp
	}

	subject := poll.Subject
	proposalsNames := make([]string, 0)
	for _, proposal := range proposals {
		proposalsNames = append(proposalsNames, proposal.Name)
	}

	err = doCreatePoll(orm, input, subject, proposalsNames)

	return err
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command.rerun",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &RerunCommand{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("command.rerun failed to build", err)
	}
}
