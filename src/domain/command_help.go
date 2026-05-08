package domain

import (
	//"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	"main/src/provider"
)

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

//func (c HelpCommand) DefineForDiscord() *disgord.ApplicationCommandOption {
//	return &disgord.ApplicationCommandOption{
//		Name:        c.GetName(),
//		Description: c.GetEmote() + " " + c.GetDescription(),
//		Type:        disgord.OptionTypeSubCommand,
//	}
//}

func (c HelpCommand) Matches(command string) bool {
	return command == HelpCommandSlug
}

func (c HelpCommand) Handle(input provider.Input) (handled bool, err error) {
	return true, handleHelpCommand(input)
}

func handleHelpCommand(input provider.Input) error {
	message := "🤖 _Hello !_ " +
		"My purpose is to help you create majority judgment polls.\n" +
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
		"> **Yes.**  Discord does not allow variadic application commands, for now, " +
		"but as a workaround you may specify multiple proposals per field, " +
		"using `|` as separator.\n" +
		"\n" +
		""

	return provider.GetResponder(input).RespondWithMessage(input, message, true)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command." + HelpCommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &HelpCommand{}
			return cmd, nil
		},
	})

	if err != nil {
		log.Fatalf("command.%s failed to build : %s\n", HelpCommandSlug, err)
	}
}
