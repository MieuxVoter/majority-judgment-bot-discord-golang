package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di/v2"
	"log"
	"log/slog"
	"main/src/container"
	db "main/src/database"
	"main/src/domain"
	"main/src/provider"
	"main/src/security"
	"main/src/services"
	"math/rand"
	"strings"
	"xorm.io/xorm"
)

const InfoCommandSlug = "info"

// InfoCommand displays miscellaneous information about the bot
type InfoCommand struct {
	orm       *xorm.Engine
	liberapay *services.Liberapay
	logger    *slog.Logger
}

func (c InfoCommand) GetTranslationKey() string {
	return "MjInfo"
}

func (c InfoCommand) GetEmote() string {
	return "🤖"
}

func (c InfoCommand) GetName() string {
	return InfoCommandSlug
}

//func (c InfoCommand) GetDescription() string {
//	return "Display miscellaneous information about me on this server"
//}

func (c InfoCommand) GetOptionsForDiscord() []discord.ApplicationCommandOption {
	return []discord.ApplicationCommandOption{}
}

func (c InfoCommand) Matches(command string) bool {
	return command == c.GetName()
}

func (c InfoCommand) Handle(input provider.Input) (err error) {
	return handleInfoCommand(c, input)
}

func handleInfoCommand(
	command InfoCommand,
	input provider.Input,
) error {
	guildVendorId, _ := input.GetGuildVendorId()
	guild, err := db.GetOrCreateGuild(command.orm, guildVendorId)
	if err != nil {
		message := "Could not access the guild.  _Suddenly, everything is on fire._ 🔥"
		return domain.RespondUserError(input, message)
	}
	allPollsAmount, errCountAll := db.CountPolls(command.orm)
	if errCountAll != nil {
		message := "Could not count the polls.  _Suddenly, Notre-Dame is on fire._ 🔥"
		return domain.RespondServerError(input, message)
	}
	guildPollsAmount, errCountGuildPolls := db.CountGuildPolls(command.orm, guild)
	if errCountGuildPolls != nil {
		message := "Could not count this guild's polls.  _Suddenly, Australia is on fire._ 🔥"
		return domain.RespondServerError(input, message)
	}
	thanksSlice := []string{
		"MieuxVoter.fr",
		"Vesporium (test)",
		"OrelSac (test)",
		"Chantal (test)",
		"Roipoussiere (code)",
		"Trollune (code)",
		//"Marjolaine Leray (illustration)",
	}
	rand.Shuffle(len(thanksSlice), func(i, j int) {
		thanksSlice[i], thanksSlice[j] = thanksSlice[j], thanksSlice[i]
	})
	thanks := strings.Join(thanksSlice, ", ")

	//survivalChance, err := command.liberapay.GetSurvivalAsString()
	//if err != nil {
	//	command.logger.Errorln("LIBERAPAY FAILED", err)
	//}

	message := "" +
		"🤖🗩 _Here is some information about myself._\n" +
		"\n" +
		"Total amount of polls by this community" + fmt.Sprintf(" : `%d`\n", guildPollsAmount) +
		"Total amount of polls across all communities" + fmt.Sprintf(" : `%d`\n", allPollsAmount) +
		"Remaining polls' quota of this community" + fmt.Sprintf(" : `%d`\n", guild.Quota) +
		"Version" + fmt.Sprintf(" : `%s`\n", security.GetVersion()) +
		"Guild Identifier" + fmt.Sprintf(" : `%s`\n", guild.Snowflake) +
		"\n" +
		"Friends" + fmt.Sprintf(" : `%s`\n", thanks) +
		"\n" +
		//"Survival Chance" + fmt.Sprintf(" : %s\n", survivalChance) +
		//"\n" +
		""

	//buttons := make([]*provider.ButtonField, 0)
	//buttons = append(buttons, &provider.ButtonField{
	//	Label: "Wish me Well",
	//	Emote: "🌠",
	//	Url:   "https://liberapay.com/MajorityJudgmentBot/",
	//})

	return provider.GetResponder(input).RespondMessage(
		input,
		message,
		true,
	)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "subcommand.mj." + InfoCommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &InfoCommand{
				orm:       ctn.Get("database.engine").(*xorm.Engine),
				liberapay: ctn.Get("liberapay").(*services.Liberapay),
				logger:    ctn.Get("logger").(*slog.Logger),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalf("subcommand.mj.%s failed to build : %s\n", InfoCommandSlug, err)
	}
}
