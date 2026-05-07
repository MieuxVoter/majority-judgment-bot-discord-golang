package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"main/src/container"
	"main/src/provider"
)

// mjSlashCommand is our main, root command /mj.
var mjSlashCommand = discord.SlashCommandCreate{
	Name:        "mj",
	Description: "Manage Majority Judgment polls",
	Options:     []discord.ApplicationCommandOption{}, // injected dynamically, see GetCommands
}

// commands are also injected dynamically, see GetCommands.
var commands = []discord.ApplicationCommandCreate{}

// areCommandsInjected marks whether the commands have been injected already.
var areCommandsInjected = false

// Command interface to implement in services declaring commands.
type Command interface {
	Define() discord.ApplicationCommandOption
	GetEmote() string
	GetName() string
	GetDescription() string
	Matches(command string) bool
	Handle(input provider.Input) (handled bool, err error)
}

// GetCommands lists all commands from services available in the container.
func GetCommands() []discord.ApplicationCommandCreate {
	if !areCommandsInjected {
		// Inject /mj subcommands from the tagged services.
		commandsServices := container.GetCollection("command.mj")
		for _, commandGeneric := range commandsServices {
			command := commandGeneric.(Command)
			fmt.Printf("registering subcommand %s\n", command.GetName())
			mjSlashCommand.Options = append(mjSlashCommand.Options, command.Define())
		}
		commands = append(commands, mjSlashCommand)

		// Inject other eventual commands later on (if any).
		// …

		// Finally, mark the commands as injected so we don't do it twice by accident.
		areCommandsInjected = true
	}

	return commands
}
