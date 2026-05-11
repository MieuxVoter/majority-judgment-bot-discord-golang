package provider

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"main/src/container"
	db "main/src/database"
	"main/src/services"
)

// ResponderInterface should be implemented by our vendor output adapters.
type ResponderInterface interface {
	Matches(input Input) bool
	RespondWithMessage(input Input, message string, ephemeral bool) error
	//RespondWithMessageAndImage(input Input, message string, imageUrl string, ephemeral bool) error
	//RespondWithMessageAndButtons(input Input, message string, buttons []*ButtonField, ephemeral bool) error
	RespondPollView(
		input Input,
		poll *db.Poll,
		proposals []*db.Proposal,
		replaceMessage bool,
	) error
	RespondWithJudgmentUi(
		input Input,
		proposal *db.Proposal,
		poll *db.Poll,
		previousJudgment *db.Judgment,
		replaceMessage bool,
	) error
	RespondBallotSummary(
		input Input,
		poll *db.Poll,
		proposals []db.Proposal,
		judgments []db.Judgment,
	) error
	RespondPollResult(
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
	//RespondBallotsInspection(
	//	input Input,
	//	poll *db.Poll,
	//	proposals []db.Proposal,
	//	judgments []db.Judgment,
	//) error
	RespondUserError(input Input, message string) error
	RespondServerError(input Input, message string) error
}

// GetResponder returns the responder adapter that matches the input provider
func GetResponder(input Input) ResponderInterface {
	responders := container.GetCollection("responder.")
	for _, genericResponder := range responders {
		responder := genericResponder.(ResponderInterface)
		if responder.Matches(input) {
			return responder
		}
	}

	// This should never happen → that's why it's fatal
	services.GetLogger().Fatalln("no matching responder found")
	return nil
}

// RaiseInvalidProviderError is sugar for raising when the provider is invalid.
// This should absolutely never raise, unless there's a terrible glitch in the code.
func RaiseInvalidProviderError(trace string) error {
	return fmt.Errorf("invalid input type for provider (%s)", trace)
}
