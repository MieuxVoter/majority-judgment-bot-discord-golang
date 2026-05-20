package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	"main/src/domain"
	"main/src/provider"
)

const ExplainCommandSlug = "explain"

type ExplainCommand struct{}

func (c ExplainCommand) GetTranslationKey() string {
	return "MjExplain"
}

func (c ExplainCommand) GetEmote() string {
	return "📖"
}

func (c ExplainCommand) GetName() string {
	return ExplainCommandSlug
}

//func (c ExplainCommand) GetDescription() string {
//	return "Explain Majority Judgment like you're five years old"
//}

func (c ExplainCommand) GetOptionsForDiscord() []discord.ApplicationCommandOption {
	return []discord.ApplicationCommandOption{}
}

func (c ExplainCommand) Matches(command string) bool {
	return command == c.GetName()
}

func (c ExplainCommand) Handle(input provider.Input) error {
	return handleExplainCommand(input)
}

func handleExplainCommand(
	input provider.Input,
) error {

	message := `
### 📚 Why Majority Judgment?

> Because it is so great that it literally poops rainbows!  :rainbow: :poop:

In common _single-choice polls_, you have to choose one _— and only one —_ proposal, and equally reject all the others, even though you probably prefer one to another amongst those rejected.

This lack of subtlety has terrible consequences on the quality of the poll and of democracy in general, such as the widely despised _bipartisan system_ in the USA and the equally hated _useful vote_ in France.

In _Majority Judgment polls_, you can give your opinion on each and every proposal.

Now you can finally say you like both bananas **and** strawberries better than spinach!  Great, isn't it?

### 🏆 How are proposals ranked in the end?

> Magic, of course!  :magic_wand:  _(just kidding)_

Proposals are ranked by their _median grade_, also nicknamed _majority grade_.

The median grade is the grade that's right in the middle ; this is why we show a vertical line in the middle of the merit profiles.

### 😸 Can I cheat?

> Not that we know of.  :sauropod: :seedling:

Majority Judgment was discovered by two scientists _(Michel Balinski & Rida Laraki)_ in 2003, and their goal was precisely to find a voting system that was _anti-strategic_.

If you think you figured out a way to cheat or to improve this bot, come over on our Discord. https://discord.gg/rAAQG9S :sparkles:
`

	//		"It is the highest given grade where at least 50% of the participants gave this grade or higher, " +
	//		"hence the _majority_ in _majority judgment_.\n" +
	//		"\n" +
	//		"Each proposal has therefore three groups of participants that emerge naturally : \n" +
	//		"the _median group_, people who gave the _median_ grade to the proposal, \n" +
	//		"the _contestation group_, people who gave a _lower_ grade than the median grade, and\n" +
	//		"the _adhesion group_, people who gave a _higher_ grade than the median grade.\n" +
	//		"\n" +
	//		"⚖ **How to sort two proposals with the same majority grade, then ?**" +
	//		"\n\n" +
	//		"We look at the _adhesion_ and _contestation_ groups of each of the two proposals, " +
	//		"and follow the decision of the biggest of these four groups.  " +
	//		"Yet again, majority prevails.\n" +
	//		""

	return domain.RespondWithMessage(input, message, true)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "subcommand.mj." + ExplainCommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &ExplainCommand{}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalf("subcommand.mj.%s failed to build : %s\n", ExplainCommandSlug, err)
	}
}
