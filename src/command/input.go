package command

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
)

// Input holds data coming from userland through the vendor.
// Trying to be generic so as to add other platforms/vendors than Discord at some point.
// This might not work, or might become troublesome, but let's strive for vendor abstraction anyway.
type Input interface {
	GetOption(command string, name string, defaultValue string) (string, error)
	GetActorVendorId() (string, error)
	GetGuildVendorId() (string, error)
	GetButtonName() (string, error)
}

// DiscordInput wrapper for data coming from Discord's userland.
type DiscordInput struct {
	Context     context.Context
	Session     disgord.Session
	Interaction *disgord.InteractionCreate
}

func (d DiscordInput) GetOption(command string, name string, defaultValue string) (string, error) {
	var options []*disgord.ApplicationCommandDataOption

	for _, option := range d.Interaction.Data.Options {
		if option.Name == command {
			options = option.Options
		}
	}

	if options == nil {
		return "", fmt.Errorf("command subquery options not found")
	}

	for _, option := range options {
		if option.Name == name {
			value := fmt.Sprintf("%s", option.Value)
			if value == "" {
				value = defaultValue
			}

			return value, nil
		}
	}

	return defaultValue, nil
}

func (d DiscordInput) GetActorVendorId() (string, error) {
	return d.Interaction.Member.UserID.String(), nil
}

func (d DiscordInput) GetGuildVendorId() (string, error) {
	return d.Interaction.GuildID.String(), nil
}

func (d DiscordInput) GetButtonName() (string, error) {
	return d.Interaction.Data.CustomID, nil
}
