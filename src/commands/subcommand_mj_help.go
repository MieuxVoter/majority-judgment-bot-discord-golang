package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	"main/src/locales"
	"main/src/provider"
)

// HelpCommandSlug is locale-insensitive (and should stay that way)
const HelpCommandSlug = "help"

type HelpCommand struct {
	localization *locales.Localization
}

func (c HelpCommand) GetTranslationKey() string {
	return "MjHelp"
}

func (c HelpCommand) GetEmote() string {
	return "👁"
}

func (c HelpCommand) GetName() string {
	return HelpCommandSlug
}

func (c HelpCommand) GetOptionsForDiscord() []discord.ApplicationCommandOption {
	return []discord.ApplicationCommandOption{}
}

func (c HelpCommand) Matches(subCommandName string) bool {
	return subCommandName == HelpCommandSlug
}

func (c HelpCommand) Handle(input provider.Input) error {
	return handleHelpCommand(input)
}

func handleHelpCommand(input provider.Input) error {
	localizer := locales.GetLocalizer(input.GetActorLanguage())
	message := ""
	message += fmt.Sprintf("🤖 _%s_ ", localizer.T("HelpHello"))
	message += localizer.T("HelpMyPurposeIs") + "\n"
	message += fmt.Sprintf("### ⚖ %s\n", localizer.T("HelpWhatIsMj"))
	message += fmt.Sprintf("> %s\n\n", localizer.T("HelpWhatIsMjAnswer"))
	message += fmt.Sprintf("### 🕵 %s\n", localizer.T("HelpCanBotReadMessage"))
	message += fmt.Sprintf("> %s\n\n", localizer.T("HelpCanBotReadMessageAnswer"))
	message += fmt.Sprintf("### ❺ %s\n", localizer.T("HelpCanUseMoreThanFiveGrades"))
	message += fmt.Sprintf("> %s\n\n", localizer.T("HelpCanUseMoreThanFiveGradesAnswer"))
	message += fmt.Sprintf("### ➏ %s\n", localizer.T("HelpCanUseMoreProposals"))
	message += fmt.Sprintf("> %s\n\n", localizer.T("HelpCanUseMoreProposalsAnswer"))

	return provider.GetResponder(input).RespondMessage(input, message, true)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "subcommand.mj." + HelpCommandSlug,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &HelpCommand{
				localization: ctn.Get("localization").(*locales.Localization),
			}
			return cmd, nil
		},
	})

	if err != nil {
		log.Fatalf("service subcommand.mj.%s failed to build : %s\n", HelpCommandSlug, err)
	}
}
