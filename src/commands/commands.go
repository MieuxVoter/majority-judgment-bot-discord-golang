package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	"main/src/provider"
)

// Command interface for the root slash command(s)
type Command interface {
	GetName() string
	GetDescription() string
}

// Subcommand interface to implement in services declaring subcommands.
type Subcommand interface {
	GetEmote() string
	GetName() string
	GetDescription() string
	Matches(subCommandName string) bool
	Handle(input provider.Input) error
	// GetOptionsForDiscord defines options for Discord, as an abstraction layer for this is work. (maybe later?)
	GetOptionsForDiscord() []discord.ApplicationCommandOption
}

// MjCommand is our main (and only) root slash command for now.
// It does not do anything by itself, and instead relies on subcommands.
type MjCommand struct{}

func (c MjCommand) GetName() string {
	return "mj"
}

func (c MjCommand) GetDescription() string {
	return "Manage Majority Judgment polls"
}

// --- Discord ---

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

	fmt.Println("Guild:", event.GuildID())

	input := provider.DiscordInput{
		Data:  data,
		Event: event,
	}
	subcommands := container.GetCollection("subcommand.mj")
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

// discordCommands are also injected dynamically, see GetDiscordCommands.
var discordCommands = []discord.ApplicationCommandCreate{}

// areDiscordCommandsInjected marks whether the commands have been injected already.
var areDiscordCommandsInjected = false

// GetDiscordCommands lists all commands from services available in the container.
func GetDiscordCommands() []discord.ApplicationCommandCreate {
	if !areDiscordCommandsInjected {
		// Inject /mj command and it subcommands from the tagged services.
		mjDiscordSlashCommand := DefineCommandForDiscord(container.Get("command.mj").(Command))
		mjSubcommandsServices := container.GetCollection("subcommand.mj")
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
		areDiscordCommandsInjected = true
	}

	return discordCommands
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command.mj",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := MjCommand{}
			return cmd, nil
		},
	})

	if err != nil {
		log.Fatalf("service command.mj failed to build: %s\n", err)
	}
}
