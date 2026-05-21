package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	db "main/src/database"
	"main/src/domain"
	"main/src/locales"
	"main/src/provider"
	"main/src/security"
	"main/src/services"
	"xorm.io/xorm"
)

// CreateMjSubcommandSlug is locale-insensitive (and should stay that way)
const CreateMjSubcommandSlug = "create"

type CreateMjSubcommand struct {
	orm          *xorm.Engine
	gradings     *services.Gradings
	localization *locales.Localization
}

func (c CreateMjSubcommand) GetTranslationKey() string {
	return "MjCreate"
}

func (c CreateMjSubcommand) GetEmote() string {
	return "➕"
}

func (c CreateMjSubcommand) GetName() string {
	return CreateMjSubcommandSlug
}

func (c CreateMjSubcommand) GetOptionsForDiscord() []discord.ApplicationCommandOption {
	return []discord.ApplicationCommandOption{
		// Note: we cannot use spaces or uppercase in Name (2026-05)
		discord.ApplicationCommandOptionString{
			Name:                     `subject`,
			Description:              `The poll's subject, such as "Meeting date"`,
			NameLocalizations:        c.localization.GetTranslations("ParameterSubject"),
			DescriptionLocalizations: c.localization.GetTranslations("ParameterSubjectDescription"),
			Required:                 true,
		},
		// *How to get variadism here, for proposals?*
		// Right now we work around the limitation with a hack,
		// by supporting adding multiple proposals per field using | as separator.
		discord.ApplicationCommandOptionString{
			Name:                     "proposal_a",
			Description:              `The name of the first proposal, like "Friday"`,
			NameLocalizations:        c.localization.GetTranslations("ParameterProposalA"),
			DescriptionLocalizations: c.localization.GetTranslations("ParameterProposalADescription"),
			Required:                 true,
		},
		discord.ApplicationCommandOptionString{
			Name:                     "proposal_b",
			Description:              `The name of the second proposal, like "Pizza"`,
			NameLocalizations:        c.localization.GetTranslations("ParameterProposalB"),
			DescriptionLocalizations: c.localization.GetTranslations("ParameterProposalBDescription"),
		},
		discord.ApplicationCommandOptionString{
			Name:                     "proposal_c",
			Description:              `The name of the third proposal, like "Beaujolais"`,
			NameLocalizations:        c.localization.GetTranslations("ParameterProposalC"),
			DescriptionLocalizations: c.localization.GetTranslations("ParameterProposalCDescription"),
		},
		discord.ApplicationCommandOptionString{
			Name:                     "proposal_d",
			Description:              `The name of the fourth proposal, like "Michel"`,
			NameLocalizations:        c.localization.GetTranslations("ParameterProposalD"),
			DescriptionLocalizations: c.localization.GetTranslations("ParameterProposalDDescription"),
		},
		discord.ApplicationCommandOptionString{
			Name:                     "proposal_e",
			Description:              `If you need more than five, use | as separator`,
			NameLocalizations:        c.localization.GetTranslations("ParameterProposalE"),
			DescriptionLocalizations: c.localization.GetTranslations("ParameterProposalEDescription"),
		},
		discord.ApplicationCommandOptionString{
			Name:        "grading",
			Description: "The grades to use in this poll",
			Choices: []discord.ApplicationCommandOptionChoiceString{
				// All the Values in here must be available as keys in [services.Gradings.Get]
				{
					Name:  "👎👍",
					Value: "👎👍",
				},
				{
					Name:  "👎🤷👍",
					Value: "👎🤷👍",
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
				// so to add more than 5 grades we need to tweak our judgment UI.
			},
		},
		//{
		//	Name:        "secrecy",
		//	Description: "Whether individual votes are kept secret or not. (default is secret)",
		//	Choices: []*discord.ApplicationCommandOptionChoiceString{
		//		{
		//			Name:  "secret for all (default)",
		//			Value: "secret",
		//		},
		//		//{
		//		//	Name:  "poll author can see",
		//		//	Value: "admin",
		//		//},
		//		{
		//			Name:  "anyone can see",
		//			Value: "public",
		//		},
		//	},
		//},
	}
}

func (c CreateMjSubcommand) Matches(command string) bool {
	return command == c.GetName()
}

func (c CreateMjSubcommand) Handle(input provider.Input) error {
	if input.IsDirectMessage() {
		message := "I can't create a poll just for you and I.  🤷  Try again in a channel with other people?"
		return provider.GetResponder(input).RespondUserError(input, message)
	}

	return handleCreateCommand(c.orm, input)
}

func handleCreateCommand(
	orm *xorm.Engine,
	input provider.Input,
) error {

	subject, err := input.GetOptionString("create", "subject", "Poll")
	if err != nil {
		return err
	}

	proposalsNames := make([]string, 0)
	for _, v := range []string{"a", "b", "c", "d", "e"} { // :(|) ooOOk?
		rawProposalName, _ := input.GetOptionString("create", "proposal_"+v, "")
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

	if len(proposalsNames) < 1 {
		return domain.RespondUserError(input, "A Poll needs at least one proposal.")
	}

	grading, _ := input.GetOptionString("create", "grading", "🤮😐😌😀🤩")
	secrecy, _ := input.GetOptionString("create", "secrecy", "secret")

	err = doCreatePoll(orm, input, subject, proposalsNames, grading, secrecy)

	return err
}

func doCreatePoll(
	orm *xorm.Engine,
	input provider.Input,
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
		return domain.RespondUserError(input, "This guild cannot create polls anymore.")
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

	err = domain.RespondPollView(input, poll, proposals, false)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "subcommand.mj." + CreateMjSubcommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &CreateMjSubcommand{
				orm:          ctn.Get("database.engine").(*xorm.Engine),
				gradings:     ctn.Get("gradings").(*services.Gradings),
				localization: ctn.Get("localization").(*locales.Localization),
			}
			return cmd, nil
		},
	})

	if err != nil {
		log.Fatalf("subcommand.mj.%s failed to build : %s\n", CreateMjSubcommandSlug, err)
	}
}
