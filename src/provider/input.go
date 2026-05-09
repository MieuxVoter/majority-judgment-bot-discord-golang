package provider

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
)

// Input holds data coming from userland through the vendor.
// Trying to be generic so as to add other platforms/vendors than Discord at some point.
// This might not work, or might become troublesome, but let's strive for vendor abstraction anyway.
type Input interface {
	GetOptionString(subcommand string, name string, defaultValue string) (string, error)
	GetActorVendorId() (string, error)
	GetActorName() (string, error)
	GetGuildVendorId() (string, error)
	IsDirectMessage() bool
}

type ButtonInput interface {
	Input
	GetButtonName() (string, error)
}

//  ___  _                   _
// |   \(_)___ __ ___ _ _ __| |
// | |) | (_-</ _/ _ \ '_/ _` |
// |___/|_/__/\__\___/_| \__,_|
// (move this to its own file?)

type DiscordInteraction interface {
	CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error
	UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error
}

// DiscordCommandInput wrapper for data coming from Discord's userland.
type DiscordCommandInput struct {
	Data  discord.SlashCommandInteractionData
	Event *handler.CommandEvent
}

func (d DiscordCommandInput) GetOptionString(subcommand string, name string, defaultValue string) (string, error) {
	option, optionWasFound := d.Data.Option(name)

	if !optionWasFound {
		return defaultValue, nil
	}

	if option.Type == discord.ApplicationCommandOptionTypeString {
		return option.String(), nil
	}

	return "", fmt.Errorf("subcommand `%s` option `%s` type unsupported", subcommand, name)
}

func (d DiscordCommandInput) GetActorVendorId() (string, error) {
	member := d.Event.Member()
	if member != nil {
		return member.User.ID.String(), nil
	}
	return "", fmt.Errorf("actor id is unavailable")
}

func (d DiscordCommandInput) GetActorName() (string, error) {
	member := d.Event.Member()
	if member != nil {
		return member.User.Username, nil
	}
	return "", fmt.Errorf("actor name is unavailable")
}

func (d DiscordCommandInput) GetGuildVendorId() (string, error) {
	guildId := d.Event.GuildID()
	if guildId != nil {
		return d.Event.GuildID().String(), nil
	}
	return "", fmt.Errorf("guild id is unavailable")
}

func (d DiscordCommandInput) GetButtonName() (string, error) {
	return d.Event.Data.CommandName(), nil
}

func (d DiscordCommandInput) IsDirectMessage() bool {
	return d.Event.GuildID() == nil
}

func (d DiscordCommandInput) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return d.Event.CreateMessage(messageCreate, opts...)
}

func (d DiscordCommandInput) UpdateMessage(_ discord.MessageUpdate, _ ...rest.RequestOpt) error {
	return fmt.Errorf("cannot update a message from a command")
}

type DiscordButtonInput struct {
	Data  discord.ButtonInteractionData
	Event *handler.ComponentEvent
}

func (d DiscordButtonInput) GetOptionString(subcommand string, name string, _ string) (string, error) {
	return "", fmt.Errorf("button does not support GetOptionString(%s, %s)", subcommand, name)
}

func (d DiscordButtonInput) GetActorVendorId() (string, error) {
	member := d.Event.Member()
	if member != nil {
		return member.User.ID.String(), nil
	}
	return "", fmt.Errorf("actor id is unavailable")
}

func (d DiscordButtonInput) GetActorName() (string, error) {
	member := d.Event.Member()
	if member != nil {
		return member.User.Username, nil
	}
	return "", fmt.Errorf("actor name is unavailable")
}

func (d DiscordButtonInput) GetGuildVendorId() (string, error) {
	guildId := d.Event.GuildID()
	if guildId != nil {
		return d.Event.GuildID().String(), nil
	}
	return "", fmt.Errorf("guild id is unavailable")
}

func (d DiscordButtonInput) GetButtonName() (string, error) {
	return d.Event.ButtonInteractionData().CustomID(), nil
}

func (d DiscordButtonInput) IsDirectMessage() bool {
	return false
}

func (d DiscordButtonInput) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return d.Event.CreateMessage(messageCreate, opts...)
}

func (d DiscordButtonInput) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return d.Event.UpdateMessage(messageUpdate, opts...)
}
