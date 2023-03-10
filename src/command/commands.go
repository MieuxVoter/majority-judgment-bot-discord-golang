package command

import (
	"context"
	"github.com/andersfylling/disgord"
)

// We'll probably want the builder pattern here instead of this static def.  Where is it?
// It would also be nice to figure out how to get variadic commands with unlimited proposals.
var commands = []*disgord.CreateApplicationCommand{
	{
		Name:        "mj",
		Description: "Manage Majority Judgment polls",
		Type:        disgord.ApplicationCommandChatInput,
		Options: []*disgord.ApplicationCommandOption{
			{
				Name:        "create",
				Description: "Create a new poll",
				Type:        disgord.OptionTypeSubCommand,
				Options: []*disgord.ApplicationCommandOption{
					{
						Name:        "subject",
						Description: "The poll's subject, such as \"When should we meet?\"",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_a",
						Description: "The name of the first proposal, like Friday",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_b",
						Description: "The name of the second proposal, like Pizza",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_c",
						Description: "The name of the third proposal, like Beaujolais",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_d",
						Description: "The name of the fourth proposal, like Michel",
						Type:        disgord.OptionTypeString,
					},
					{
						Name:        "proposal_e",
						Description: "The name of the fifth element, like Moultipass",
						Type:        disgord.OptionTypeString,
					},
					// /!. Discord limits messages integrations to 5 action rows,
					//     so we'd need multiple messages to handle more than 5 proposals.
					//     No point in adding proposal_f here for now, it won't work as-is.
					// > Well, now we use one message per proposal, but how to get variadism here?
				},
			},
			{
				Name:        "help",
				Description: "Send an SOS: ... --- ...",
				Type:        disgord.OptionTypeSubCommand,
			},
		},
	},
}

type Input struct {
	Context     context.Context
	Session     disgord.Session
	Interaction *disgord.InteractionCreate
}

type Command interface {
	Define() *disgord.ApplicationCommandOption
	Matches(command string) bool
	Handle(input *Input) (handled bool, err error)
}

func GetCommands() []*disgord.CreateApplicationCommand {
	return commands
}

func init() {
	//fmt.Println("init commands")
}
