package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	"main/src/domain"
	"main/src/locales"
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
	localizer := locales.GetLocalizer(input.GetActorLanguage())
	message := ""
	message += fmt.Sprintf("### 📚 %s\n", localizer.T("ExplainWhyMj"))
	message += fmt.Sprintf("> %s\n\n", localizer.T("ExplainWhyMjSubtitle"))
	message += fmt.Sprintf("%s\n", localizer.T("ExplainWhyMjParagraph"))
	message += fmt.Sprintf("### 🏆 %s\n", localizer.T("ExplainHowMj"))
	message += fmt.Sprintf("> %s\n\n", localizer.T("ExplainHowMjSubtitle"))
	message += fmt.Sprintf("%s\n", localizer.T("ExplainHowMjParagraph"))
	message += fmt.Sprintf("### 😸 %s\n", localizer.T("ExplainCheatMj"))
	message += fmt.Sprintf("> %s\n\n", localizer.T("ExplainCheatMjSubtitle"))
	message += fmt.Sprintf("%s  ", localizer.T("ExplainCheatMjParagraph"))
	message += "https://discord.gg/k9YRuZPSZs :sparkles:"

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
	//		"Yet again, majority prevails."

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
