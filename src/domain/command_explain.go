package domain

import (
	//"github.com/andersfylling/disgord"
	"github.com/sarulabs/di"
	"log"
	"main/src/container"
	"main/src/provider"
)

const ExplainCommandSlug = "explain"

type ExplainCommand struct{}

func (c ExplainCommand) GetEmote() string {
	return "📖"
}

func (c ExplainCommand) GetName() string {
	return ExplainCommandSlug
}

func (c ExplainCommand) GetDescription() string {
	return "Explain Majority Judgment like you're five years old"
}

//func (c ExplainCommand) Define() *disgord.ApplicationCommandOption {
//	return &disgord.ApplicationCommandOption{
//		Name:        c.GetName(),
//		Description: c.GetEmote() + " " + c.GetDescription(),
//		Type:        disgord.OptionTypeSubCommand,
//	}
//}

func (c ExplainCommand) Matches(command string) bool {
	return command == c.GetName()
}

func (c ExplainCommand) Handle(input provider.Input) (handled bool, err error) {
	return true, handleExplainCommand(input)
}

func handleExplainCommand(input provider.Input) error {
	message := "📚 **Majority Judgment, Why ?**" +
		"\n\n" +
		"In common _uninominal polls_, participants each choose one (and only one) proposal, " +
		"and equally reject all the others, even though they may prefer one to another.\n" +
		"This lack of subtlety yields appalling consequences for the health of the democratic discourse, " +
		"such as the widely despised _bipartisan_ system in the USA and the so-called _useful vote_.\n" +
		"\n" +
		"In _majority judgment polls_, participants may judge each proposal individually.\n" +
		"\n" +
		"🏆 **How are proposals ranked in the end ?**" +
		"\n\n" +
		"Proposals are ranked by their _median grade_, also nicknamed _majority grade_.\n" +
		"It is the highest given grade where at least 50% of the participants gave this grade or higher, " +
		"hence the _majority_ in _majority judgment_.\n" +
		"\n" +
		"Each proposal has therefore three groups of participants that emerge naturally : \n" +
		"the _median group_, people who gave the _median_ grade to the proposal, \n" +
		"the _contestation group_, people who gave a _lower_ grade than the median grade, and\n" +
		"the _adhesion group_, people who gave a _higher_ grade than the median grade.\n" +
		"\n" +
		"⚖ **How to sort two proposals with the same majority grade, then ?**" +
		"\n\n" +
		"We look at the _adhesion_ and _contestation_ groups of each of the two proposals, " +
		"and follow the decision of the biggest of these four groups.  " +
		"Yet again, majority prevails.\n" +
		"\n" +
		"\n" +
		"_(illustration: MarjolaineLeray.com)_\n" +
		""
	// FIXME: find a another space for this image somewhere on the internet  (permalink, please)
	imageUrl := "https://media.discordapp.net/attachments/855665583869919233/1087985229177831475/equality_explained.png?width=652&height=493"

	return provider.GetResponder(input).RespondWithMessageAndImage(input, message, imageUrl, true)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command." + ExplainCommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &ExplainCommand{}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalf("command.%s failed to build : %s\n", ExplainCommandSlug, err)
	}
}
