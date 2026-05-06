package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var test = discord.SlashCommandCreate{
	Name:        "test",
	Description: "test command",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "choice",
			Description:  "some autocomplete choice",
			Required:     true,
			Autocomplete: true,
		},
	},
}

func TestHandler(e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{}.
		WithContentf("I AM A DISEMBODIED HEAD SO YOU CAN TRUST ME ; YOU CAN TRUST ME"),
	)
	//return e.CreateMessage(discord.NewMessageCreateBuilder().
	//	SetContentf("test command. Choice: %s", e.SlashCommandInteractionData().String("choice")).
	//	AddActionRow(discord.NewPrimaryButton("test", "/test-button")).
	//	Build(),
	//)
}
