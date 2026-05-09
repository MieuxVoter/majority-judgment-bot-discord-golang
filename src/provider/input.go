package provider

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Input holds data coming from userland through the vendor.
// Trying to be generic so as to add other platforms/vendors than Discord at some point.
// This might not work, or might become troublesome, but let's strive for vendor abstraction anyway.
type Input interface {
	GetOptionString(subcommand string, name string, defaultValue string) (string, error)
	GetActorVendorId() (string, error)
	GetActorName() (string, error)
	GetGuildVendorId() (string, error)
	GetButtonName() (string, error)
	IsDirectMessage() bool
}

//  ___  _                   _
// |   \(_)___ __ ___ _ _ __| |
// | |) | (_-</ _/ _ \ '_/ _` |
// |___/|_/__/\__\___/_| \__,_|
// (move this to its own file plz)

// DiscordInput wrapper for data coming from Discord's userland.
type DiscordInput struct {
	Data  discord.SlashCommandInteractionData
	Event *handler.CommandEvent
}

func (d DiscordInput) GetOptionString(subcommand string, name string, defaultValue string) (string, error) {
	option, optionWasFound := d.Data.Option(name)

	if !optionWasFound {
		return defaultValue, nil
	}

	if option.Type == discord.ApplicationCommandOptionTypeString {
		return option.String(), nil
	}

	return "", fmt.Errorf("subcommand `%s` option `%s` type unsupported", subcommand, name)
}

func (d DiscordInput) GetActorVendorId() (string, error) {
	member := d.Event.Member()
	if member != nil {
		return member.User.ID.String(), nil
	}
	return "", fmt.Errorf("actor id is unavailable")
}

func (d DiscordInput) GetActorName() (string, error) {
	member := d.Event.Member()
	if member != nil {
		return member.User.Username, nil
	}
	return "", fmt.Errorf("actor name is unavailable")
}

func (d DiscordInput) GetGuildVendorId() (string, error) {
	guildId := d.Event.GuildID()
	if guildId != nil {
		return d.Event.GuildID().String(), nil
	}
	return "", fmt.Errorf("guild id is unavailable")
}

func (d DiscordInput) GetButtonName() (string, error) {
	return d.Event.Data.CommandName(), nil
}

func (d DiscordInput) IsDirectMessage() bool {
	return d.Event.GuildID() == nil
}
