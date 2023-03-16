package domain

import (
	"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/security"
	"xorm.io/xorm"
)

type CreateCommand struct {
	orm *xorm.Engine
}

func (c CreateCommand) Define() *disgord.ApplicationCommandOption {
	return &disgord.ApplicationCommandOption{
		Name:        "create",
		Description: "Create a new poll",
		Type:        disgord.OptionTypeSubCommand,
		Options: []*disgord.ApplicationCommandOption{
			{
				Type:        disgord.OptionTypeString,
				Name:        "subject",
				Description: "The poll's subject, such as \"When should we meet?\"",
			},
			// How to get variadism here, for proposals?
			{
				Type:        disgord.OptionTypeString,
				Name:        "proposal_a",
				Description: "The name of the first proposal, like Friday",
			},
			{
				Type:        disgord.OptionTypeString,
				Name:        "proposal_b",
				Description: "The name of the second proposal, like Pizza",
			},
			{
				Type:        disgord.OptionTypeString,
				Name:        "proposal_c",
				Description: "The name of the third proposal, like Beaujolais",
			},
			{
				Type:        disgord.OptionTypeString,
				Name:        "proposal_d",
				Description: "The name of the fourth proposal, like Michel",
			},
			{
				Type:        disgord.OptionTypeString,
				Name:        "proposal_e",
				Description: "The name of the fifth element, like Moultipass",
			},
			{
				Type:        disgord.OptionTypeString,
				Name:        "grading",
				Description: "The grades to use in this poll",
				Choices: []*disgord.ApplicationCommandOptionChoice{
					{
						Name:  "👎👍",
						Value: "👎👍",
					},
					{
						Name:  "👎👊👍",
						Value: "👎👊👍",
					},
					{
						Name:  "🤮😐😀🤩",
						Value: "🤮😐😀🤩",
					},
					{
						Name:  "🤮😐😌😀🤩 (default)",
						Value: "🤮😐😌😀🤩",
					},
					// Discord only supports at most 5 buttons per action row,
					// so to add more of those we need to tweak our judgment UI.
				},
			},
			{
				Type:        disgord.OptionTypeString,
				Name:        "secrecy",
				Description: "Whether individual votes are kept secret or not. (default is secret)",
				Choices: []*disgord.ApplicationCommandOptionChoice{
					{
						Name:  "secret for all (default)",
						Value: "secret",
					},
					//{
					//	Name:  "poll author can see",
					//	Value: "admin",
					//},
					{
						Name:  "anyone can see",
						Value: "public",
					},
				},
			},
		},
	}
}

func (c CreateCommand) Matches(command string) bool {
	return command == "create"
}

func (c CreateCommand) Handle(input Input) (handled bool, err error) {
	if input.IsDirectMessage() {
		message := "I can't create a poll just for you and I.  🤷  Try again in a channel with other people?"
		return true, RespondUserError(input, message)
	}

	return true, handleCreateCommand(c.orm, input)
}

func handleCreateCommand(
	orm *xorm.Engine,
	input Input,
) error {

	subject, err := input.GetOption("create", "subject", "Poll")
	proposalsNames := make([]string, 0)
	for _, v := range []string{"a", "b", "c", "d", "e"} { // :(|) ooOOk?
		rawProposalName, _ := input.GetOption("create", "proposal_"+v, "")
		if rawProposalName == "" {
			continue
		}
		// Discord does not accept variadic commands yet, so we're accepting multiple proposals
		// in each of the proposal_x fields, using the character | as separator.
		// To use the | character in your proposal names, double it.
		compoundProposalsNames := security.ExtractProposalsNames(rawProposalName)

		for _, proposalName := range compoundProposalsNames {
			proposalsNames = append(proposalsNames, security.RemoveMarkdown(proposalName))
		}
	}

	if len(proposalsNames) < 2 {
		err = RespondUserError(input, "A Poll needs at least two proposals.")
		if err != nil {
			return err
		}
		return nil
	}

	grading, _ := input.GetOption("create", "grading", "🤮😐😌😀🤩")
	secrecy, _ := input.GetOption("create", "secrecy", "secret")

	err = doCreatePoll(orm, input, subject, proposalsNames, grading, secrecy)

	return err
}

func doCreatePoll(
	orm *xorm.Engine,
	input Input,
	subject string,
	proposalsNames []string,
	grading string,
	secrecy string,
) error {

	guildVendorId, _ := input.GetGuildVendorId()
	guild, err := db.GetOrCreateGuild(orm, guildVendorId)
	if err != nil {
		return err
	}

	// Check if the guild is allowed to create new polls
	isAllowed, err := security.CanGuildCreatePoll(orm, guild)
	if err != nil {
		return err
	}
	if !isAllowed {
		err = RespondUserError(input, "This guild cannot create polls anymore.")
		if err != nil {
			return err
		}
		return nil
	}

	// Decrement the guild's quota
	if guild.Quota > 0 {
		guild.Quota = guild.Quota - 1
	}
	_, err = orm.
		Cols("quota").
		Where("snowflake = ?", guild.Snowflake).
		Update(guild)
	if err != nil {
		return err
	}

	poll := &db.Poll{
		Subject: subject,
		GuildId: guild.Id,
		Grading: grading,
		Secrecy: secrecy,
	}
	_, err = orm.InsertOne(poll)
	if err != nil {
		return err
	}

	proposals := make([]*db.Proposal, 0)
	for _, proposalName := range proposalsNames {
		proposal := &db.Proposal{
			Name:   proposalName,
			PollId: poll.Id,
		}
		proposals = append(proposals, proposal)
	}
	_, err = orm.Insert(&proposals)
	if err != nil {
		return err
	}

	err = RespondWithPollUi(input, poll, proposals, false)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command.create",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &CreateCommand{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("command.create failed to build", err)
	}
}
