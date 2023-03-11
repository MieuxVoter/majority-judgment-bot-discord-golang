package command

import (
	"context"
	"github.com/andersfylling/disgord"
	"main/src/container"
)

// We'll probably want the builder pattern here instead of this static def.  Where is it?
// It would also be nice to figure out how to get variadic commands, for unlimited proposals.
var commands = []*disgord.CreateApplicationCommand{
	{
		Name:        "mj",
		Description: "Manage Majority Judgment polls",
		Type:        disgord.ApplicationCommandChatInput,
		Options:     []*disgord.ApplicationCommandOption{}, // injected dynamically, see GetCommands
	},
}
var areCommandsDefined = false

// Input wrapper for data coming from userland
type Input struct {
	Context     context.Context
	Session     disgord.Session
	Interaction *disgord.InteractionCreate
}

// Command interface to implement in services declaring commands
type Command interface {
	Define() *disgord.ApplicationCommandOption
	Matches(command string) bool
	Handle(input *Input) (handled bool, err error)
}

func GetCommands() []*disgord.CreateApplicationCommand {
	if !areCommandsDefined {
		commandsServices := container.GetCollection("command")
		for _, commandGeneric := range commandsServices {
			command := commandGeneric.(Command)
			commands[0].Options = append(commands[0].Options, command.Define())
		}
		areCommandsDefined = true
	}

	return commands
}

func init() {
	//fmt.Println("init commands")
}
