package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"main/src/container"
	"main/src/provider"
)

func DefineCommandForDiscord(c Command) discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        c.GetName(),
		Description: c.GetDescription(),
		Options:     []discord.ApplicationCommandOption{}, // injected dynamically, see GetDiscordCommands
	}
}

func DefineSubcommandForDiscord(sc Subcommand) discord.ApplicationCommandOption {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        sc.GetName(),
		Description: sc.GetEmote() + " " + sc.GetDescription(),
		Options:     sc.GetOptionsForDiscord(),
	}
}

func MjDiscordSlashCommandHandler(data discord.SlashCommandInteractionData, event *handler.CommandEvent) error {
	if data.SubCommandName == nil {
		// Note: I have not found any way to trigger this situation yet.
		return event.CreateMessage(discord.MessageCreate{}.
			WithContentf(":party: **ACHIEVEMENT UNLOCKED**: _Nifty Haxxor_ :party:"),
		)
	}

	//fmt.Println("Guild:", event.GuildID())

	input := provider.DiscordCommandInput{
		Data:  data,
		Event: event,
	}
	subcommands := container.GetCollection("subcommand.mj.")
	for _, subcommand := range subcommands {
		if !subcommand.(Subcommand).Matches(*data.SubCommandName) {
			continue
		}

		return subcommand.(Subcommand).Handle(input)
	}

	// Note: I have not found any way to trigger this situation yet.
	return event.CreateMessage(discord.MessageCreate{}.
		WithContentf(":party: **ACHIEVEMENT UNLOCKED**: _404: Hack Not Found_ :party:"),
	)
}

// discordCommands are also collected and injected in here dynamically, see GetDiscordCommands.
var discordCommands []discord.ApplicationCommandCreate

// areDiscordCommandsCollected marks whether the commands have been collected already.
// We could probably dispense of this boolean and instead check whether discordCommands is empty.
var areDiscordCommandsCollected = false

// GetDiscordCommands lists all commands from services available in the container.
func GetDiscordCommands() []discord.ApplicationCommandCreate {
	if !areDiscordCommandsCollected {
		// Inject /mj command and it subcommands from the tagged services.
		mjDiscordSlashCommand := DefineCommandForDiscord(container.Get("command.mj").(Command))
		mjSubcommandsServices := container.GetCollection("subcommand.mj.")
		for _, subcommandGeneric := range mjSubcommandsServices {
			subcommand := subcommandGeneric.(Subcommand)
			mjDiscordSlashCommand.Options = append(
				mjDiscordSlashCommand.Options,
				DefineSubcommandForDiscord(subcommand),
			)
		}
		discordCommands = append(discordCommands, mjDiscordSlashCommand)

		// Inject other root Discord commands later on (if any; none is envisioned).
		// …

		// Finally, toggle our memoization marker.
		areDiscordCommandsCollected = true
	}

	return discordCommands
}
