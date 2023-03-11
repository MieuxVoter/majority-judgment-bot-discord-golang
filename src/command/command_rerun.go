package command

import (
	"context"
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

func (c RerunCommand) Handle(input *Input) (handled bool, err error) {
	return true, handleRerunCommand(
		c.orm,
		input.Context,
		input.Session,
		input.Interaction,
	)
}

func handleRerunCommand(
	orm *xorm.Engine,
	c context.Context,
	s disgord.Session,
	h *disgord.InteractionCreate,
) error {

	subcommandOptions, err := getSubcommandOptions(h.Data.Options, "rerun")
	if err != nil {
		return err
	}

	//subject := getOptionStringByName(subcommandOptions, "subject", "Poll")
	//proposalsNames := make([]string, 0)
	//for _, v := range []string{"a", "b", "c", "d", "e"} { // :(|) ooOOk?
	//	proposalName := getOptionStringByName(subcommandOptions, "proposal_"+v, "")
	//	if proposalName == "" {
	//		continue
	//	}
	//	proposalsNames = append(proposalsNames, proposalName)
	//}
	//if len(proposalsNames) < 2 {
	//	err = RespondCommandFailure(c, s, h, "A Poll needs at least two proposals.")
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//}

	pollIdString := getOptionStringByName(subcommandOptions, "poll", "")
	if pollIdString == "" {
		// TODO: fetch most recent poll on this guild
		return RespondCommandFailure(c, s, h, "Fetching most recent poll is not implemented yet.  "+
			"Please provide a poll identifier.")
	}
	pollIdString = strings.Trim(pollIdString, "#")

	pollId, errConv := strconv.Atoi(pollIdString)
	if errConv != nil {
		return errConv
	}
	poll, errPoll := db.GetPoll(orm, uint64(pollId))
	if errPoll != nil {
		return errPoll
	}
	if poll == nil {
		message := "The specified poll was not found.  " +
			"Sorry.  _Here, have a banana instead : 🍌_"
		return RespondCommandFailure(c, s, h, message)
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

	err = doCreatePoll(orm, c, s, h, subject, proposalsNames)

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
