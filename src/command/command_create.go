package command

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
				Name:        "subject",
				Description: "The poll's subject, such as \"When should we meet?\"",
				Type:        disgord.OptionTypeString,
			},
			{
				Name:        "proposal_a",
				Description: "The name of the first proposal, like Friday",
				Type:        disgord.OptionTypeString,
			},
			{
				Name:        "proposal_b",
				Description: "The name of the second proposal, like Pizza",
				Type:        disgord.OptionTypeString,
			},
			{
				Name:        "proposal_c",
				Description: "The name of the third proposal, like Beaujolais",
				Type:        disgord.OptionTypeString,
			},
			{
				Name:        "proposal_d",
				Description: "The name of the fourth proposal, like Michel",
				Type:        disgord.OptionTypeString,
			},
			{
				Name:        "proposal_e",
				Description: "The name of the fifth element, like Moultipass",
				Type:        disgord.OptionTypeString,
			},
			// /!. Discord limits messages integrations to 5 action rows,
			//     so we'd need multiple messages to handle more than 5 proposals.
			//     No point in adding proposal_f here for now, it won't work as-is.
			// > Well, now we use one message per proposal, but how to get variadism here?
		},
	}
}

func (c CreateCommand) Matches(command string) bool {
	return command == "create"
}

func (c CreateCommand) Handle(input Input) (handled bool, err error) {
	return true, handleCreateCommand(
		c.orm,
		input,
	)
}

// handleCreateCommand is to refactor into a class at some point
func handleCreateCommand(
	x *xorm.Engine,
	input Input,
) error {

	subject, err := input.GetOption("create", "subject", "Poll")
	proposalsNames := make([]string, 0)
	for _, v := range []string{"a", "b", "c", "d", "e"} { // :(|) ooOOk?
		proposalName, _ := input.GetOption("create", "proposal_"+v, "")
		if proposalName == "" {
			continue
		}
		proposalsNames = append(proposalsNames, proposalName)
	}

	if len(proposalsNames) < 2 {
		err = RespondUserError(input, "A Poll needs at least two proposals.")
		if err != nil {
			return err
		}
		return nil
	}

	err = doCreatePoll(x, input, subject, proposalsNames)

	return err
}

func doCreatePoll(
	orm *xorm.Engine,
	input Input,
	subject string,
	proposalsNames []string,
) error {

	guildVendorId, _ := input.GetGuildVendorId()
	guild, err := db.GetGuild(orm, guildVendorId)
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
