package commands

import (
	"github.com/disgoorg/disgo/discord"
	"main/src/provider"
)

// Command interface for the root slash command(s)
type Command interface {
	GetName() string
	GetDescription() string
}

// Subcommand interface to implement in services declaring subcommands.
type Subcommand interface {
	GetTranslationKey() string
	GetEmote() string
	GetName() string
	GetDescription() string
	Matches(subCommandName string) bool
	Handle(input provider.Input) error
	// GetOptionsForDiscord defines options for Discord, as an abstraction layer for this is work. (maybe later?)
	GetOptionsForDiscord() []discord.ApplicationCommandOption
}
