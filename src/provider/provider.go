package provider

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"main/src/container"
	db "main/src/database"
	"main/src/services"
)

// ResponderInterface should be implemented by our vendor output adapters
type ResponderInterface interface {
	Matches(input Input) bool
	RespondWithMessage(input Input, message string, ephemeral bool) error
	RespondWithMessageAndImage(input Input, message string, imageUrl string, ephemeral bool) error
	RespondWithJudgmentUi(
		input Input,
		proposal *db.Proposal,
		poll *db.Poll,
		previousJudgment *db.Judgment,
		replaceMessage bool,
	) error
	RespondPollView(
		input Input,
		poll *db.Poll,
		proposals []*db.Proposal,
		replaceMessage bool,
	) error
	RespondDeliberation(
		input Input,
		poll *db.Poll,
		proposals []db.Proposal,
		pollTally *judgment.PollTally,
		pollResult *judgment.PollResult,
		title string,
		message string,
		asPrivateMessage bool,
		canInspect bool,
	) error
	RespondBallotsInspection(
		input Input,
		poll *db.Poll,
		proposals []db.Proposal,
		judgments []db.Judgment,
	) error
	RespondUserError(input Input, message string) error
	RespondServerError(input Input, message string) error
}

// GetResponder returns the responder adapter that matches the input provider
func GetResponder(input Input) ResponderInterface {
	responders := container.GetCollection("responder")
	for _, genericResponder := range responders {
		responder := genericResponder.(ResponderInterface)
		if responder.Matches(input) {
			return responder
		}
	}

	services.GetLogger().Fatalln("no matching responder found")
	return nil
}

func RaiseInvalidProviderError(trace string) error {
	return fmt.Errorf("invalid input type for provider (%s)", trace)
}
