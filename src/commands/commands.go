package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"main/src/container"
	"main/src/provider"
)

// mjSlashCommand is our main, root command /mj.
var mjSlashCommand = discord.SlashCommandCreate{
	Name:        "mj",
	Description: "Manage Majority Judgment polls",
	Options:     []discord.ApplicationCommandOption{}, // injected dynamically, see GetCommands
}

func MjSlashCommandHandler(data discord.SlashCommandInteractionData, event *handler.CommandEvent) error {
	//fmt.Println("MjSlashCommandHandler called")
	//err := e.CreateMessage(discord.MessageCreate{}.
	//	WithContentf("HELP ME I AM TRAPPED IN A BOT FACTORY"),
	//)

	if data.SubCommandName == nil {
		return event.CreateMessage(discord.MessageCreate{}.
			WithContentf(":party: **ACHIEVEMENT UNLOCKED**: _Nifty Haxxor_ :party:"),
		)
	}

	input := provider.DiscordInput{
		Data:  data,
		Event: event,
	}
	subCommands := container.GetCollection("command.mj")
	for _, subCommand := range subCommands {
		if !subCommand.(Command).Matches(*data.SubCommandName) {
			continue
		}

		return subCommand.(Command).Handle(input)
	}

	return event.CreateMessage(discord.MessageCreate{}.
		WithContentf(":party: **ACHIEVEMENT UNLOCKED**: _404: Hack Not Found_ :party:"),
	)
}

// commands are also injected dynamically, see GetCommands.
var commands = []discord.ApplicationCommandCreate{}

// areCommandsInjected marks whether the commands have been injected already.
var areCommandsInjected = false

// Command interface to implement in services declaring commands.
type Command interface {
	GetEmote() string
	GetName() string
	GetDescription() string
	Matches(subCommandName string) bool
	Handle(input provider.Input) error
	DefineForDiscord() discord.ApplicationCommandOption
}

// GetCommands lists all commands from services available in the container.
func GetCommands() []discord.ApplicationCommandCreate {
	if !areCommandsInjected {
		// Inject /mj command and it subcommands from the tagged services.
		commandsServices := container.GetCollection("command.mj")
		for _, commandGeneric := range commandsServices {
			command := commandGeneric.(Command)
			//fmt.Printf("registering subcommand %s\n", command.GetName())
			mjSlashCommand.Options = append(mjSlashCommand.Options, command.DefineForDiscord())
		}
		commands = append(commands, mjSlashCommand)

		// Inject other root commands later on (if any; none is envisioned).
		// …

		// Finally, mark the commands as injected so we don't do it twice by accident.
		areCommandsInjected = true
	}

	return commands
}
