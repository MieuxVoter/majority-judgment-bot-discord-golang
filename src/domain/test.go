package domain

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
	//choice := e.SlashCommandInteractionData().String("choice")
	return e.CreateMessage(discord.MessageCreate{}.
		WithContentf("I AM A DISEMBODIED HEAD SO YOU CAN TRUST ME … YOU CAN TRUST ME"),
	)
}

func TestAutocompleteHandler(e *handler.AutocompleteEvent) error {
	return e.AutocompleteResult([]discord.AutocompleteChoice{
		discord.AutocompleteChoiceString{
			Name:  "Dominique",
			Value: "1",
		},
		discord.AutocompleteChoiceString{
			Name:  "Giulia",
			Value: "2",
		},
		discord.AutocompleteChoiceString{
			Name:  "Hayley",
			Value: "3",
		},
	})
}
