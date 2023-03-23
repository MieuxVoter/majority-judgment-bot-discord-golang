package provider

import "github.com/andersfylling/disgord"

type ButtonField struct {
	Id    string
	Style disgord.ButtonStyle
	Label string
	Emote string
	Url   string
}
