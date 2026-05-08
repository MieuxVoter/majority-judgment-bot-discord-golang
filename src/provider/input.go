package provider

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Input holds data coming from userland through the vendor.
// Trying to be generic so as to add other platforms/vendors than Discord at some point.
// This might not work, or might become troublesome, but let's strive for vendor abstraction anyway.
type Input interface {
	GetOption(subcommand string, name string, defaultValue string) (string, error)
	//GetActorVendorId() (string, error)
	//GetActorName() (string, error)
	GetGuildVendorId() (string, error)
	//GetButtonName() (string, error)
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

// FIXME
func (d DiscordInput) GetOption(subcommand string, name string, defaultValue string) (string, error) {
	//var options []*discord.ApplicationCommandDataOption
	//
	//for _, option := range d.Interaction.Data.Options {
	//	if option.Name == subcommand {
	//		options = option.Options
	//	}
	//}
	//
	//if options == nil {
	//	return "", fmt.Errorf("command subquery options not found")
	//}
	//
	//for _, option := range options {
	//	if option.Name == name {
	//		value := fmt.Sprintf("%s", option.Value)
	//		if value == "" {
	//			value = defaultValue
	//		}
	//
	//		return value, nil
	//	}
	//}

	return defaultValue, nil
}

//func (d DiscordInput) GetActorVendorId() (string, error) {
//	return d.Interaction.Member.UserID.String(), nil
//}

//func (d DiscordInput) GetActorName() (string, error) {
//	return d.Interaction.Member.User.Username, nil
//}

func (d DiscordInput) GetGuildVendorId() (string, error) {
	return d.Event.GuildID().String(), nil
}

//func (d DiscordInput) GetButtonName() (string, error) {
//	return d.Interaction.Data.CustomID, nil
//}

func (d DiscordInput) IsDirectMessage() bool {
	return false // FIXME
	//return d.Event.UserCommandInteractionData().GuildID()
	//return d.Interaction.GuildID.IsZero()
}
