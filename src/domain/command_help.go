package domain

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
)

type HelpCommand struct{}

func (c HelpCommand) Define() *disgord.ApplicationCommandOption {
	return &disgord.ApplicationCommandOption{
		Name:        "help",
		Description: "Send an SOS: ... --- ...",
		Type:        disgord.OptionTypeSubCommand,
	}
}

func (c HelpCommand) Matches(command string) bool {
	return command == "help"
}

func (c HelpCommand) Handle(input Input) (handled bool, err error) {
	if d, ok := (input).(DiscordInput); ok {
		return true, handleHelpCommand(d)
	}
	return false, fmt.Errorf("unknown vendor")
}

func handleHelpCommand(
	input DiscordInput,
) error {
	err := input.Session.SendInteractionResponse(input.Context, input.Interaction, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Flags: disgord.MessageFlagEphemeral,
			Content: "🤖 _Hello !_ " +
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
				"but as a workaround may specify multiple proposals per field, using `|` as separator.\n" +
				"\n" +
				"",
		},
	})

	return err
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command.help",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &HelpCommand{}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("command.help failed to build", err)
	}
}
