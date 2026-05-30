package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"main/src/container"
	"main/src/locales"
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
	dsc := discord.ApplicationCommandOptionSubCommand{
		Name:    sc.GetName(), // must be constant; used to detect subcommand when handling
		Options: sc.GetOptionsForDiscord(),
	}

	localization := locales.GetLocalization()
	nameLocalizations := make(map[discord.Locale]string)
	descriptionLocalizations := make(map[discord.Locale]string)

	for languageIndex, language := range localization.GetLanguages() {
		localizer := localization.GetLocalizer(language)
		locale := discord.Locale(language)

		nameLocalization := localizer.T(
			fmt.Sprintf("Command%sName", sc.GetTranslationKey()),
		)
		descriptionLocalization := localizer.T(
			fmt.Sprintf("Command%sDescription", sc.GetTranslationKey()),
		)
		if languageIndex == 0 { // the default language is always the first in the list
			//dsc.Name = nameLocalization
			dsc.Description = descriptionLocalization
		}
		if nameLocalization != "" {
			nameLocalizations[locale] = nameLocalization
		}
		if descriptionLocalization != "" {
			descriptionLocalizations[locale] = sc.GetEmote() + " " + descriptionLocalization
		}
	}

	dsc.NameLocalizations = nameLocalizations
	dsc.DescriptionLocalizations = descriptionLocalizations

	return dsc
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

// discordCommands are collected and memoized in here dynamically, see GetDiscordCommands.
var discordCommands []discord.ApplicationCommandCreate

// areDiscordCommandsCollected marks whether the commands have been collected already.
// We could probably dispense of this boolean and instead check whether discordCommands is empty?
var areDiscordCommandsCollected = false

// GetDiscordCommands lists all commands from services available in the container.
func GetDiscordCommands() []discord.ApplicationCommandCreate {
	if !areDiscordCommandsCollected {
		// Collect /mj command and it subcommands from the tagged services.
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

		// Collect other root Discord commands (if any; none is envisioned).
		// …

		// Finally, toggle our memoization marker.
		areDiscordCommandsCollected = true
	}

	return discordCommands
}
