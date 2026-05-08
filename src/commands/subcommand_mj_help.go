package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	"main/src/provider"
)

// HelpCommandSlug is locale-insensitive (and should stay that way)
const HelpCommandSlug = "help"

type HelpCommand struct{}

func (c HelpCommand) GetEmote() string {
	return "👁"
}

func (c HelpCommand) GetName() string {
	return HelpCommandSlug
}

func (c HelpCommand) GetDescription() string {
	return "General help about how to interact with me"
}

func (c HelpCommand) GetOptionsForDiscord() []discord.ApplicationCommandOption {
	return []discord.ApplicationCommandOption{}
}

func (c HelpCommand) Matches(subCommandName string) bool {
	return subCommandName == HelpCommandSlug
}

func (c HelpCommand) Handle(input provider.Input) error {
	return handleHelpCommand(input)
}

func handleHelpCommand(input provider.Input) error {
	message := "🤖 _Hello !_ " +
		"My purpose is to help you create Majority Judgment polls.\n" +
		"\n" +
		"Try me out:\n" +
		"⌨ `/mj create <subject> <proposal_a> <proposal_b> …`\n" +
		"\n" +
		"⚖ **What is Majority Judgment?**\n" +
		"> A pretty rad **voting system.**  It is used in french :flag_fr: wine 🍷 contests. " +
		"It is simple, subtle and fair.\n" +
		"\n" +
		"🕵 **Can this bot read our messages?**\n" +
		"> **No.**  For extra privacy, this modern bot is NOT allowed to read messages, " +
		"only react to its own `/mj` command and button interactions.\n" +
		"\n" +
		"❺ **Can I use more than 5 grades?**\n" +
		"> **Not for now.**  Discord limits messages to 5 buttons per action row, " +
		"so we'll need more wit to support more grades.\n" +
		"\n" +
		"❺ **Can I use more than 5 proposals?**\n" +
		"> **Yes.**  Discord does not allow variadic application discordCommands, for now, " +
		"but as a workaround you may specify multiple proposals per field, " +
		"using `|` as separator.\n" +
		"\n" +
		""

	return provider.GetResponder(input).RespondWithMessage(input, message, true)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "subcommand.mj." + HelpCommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &HelpCommand{}
			return cmd, nil
		},
	})

	if err != nil {
		log.Fatalf("service subcommand.mj.%s failed to build : %s\n", HelpCommandSlug, err)
	}
}
