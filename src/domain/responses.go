package domain

import (
	db "main/src/database"
	"main/src/provider"
)

func RespondPollView(
	input provider.Input,
	poll *db.Poll,
	proposals []*db.Proposal,
	replaceMessage bool,
) error {
	return provider.GetResponder(input).RespondPollView(
		input,
		poll,
		proposals,
		replaceMessage,
	)
}

func RespondWithJudgmentUi(
	input provider.Input,
	proposal *db.Proposal,
	poll *db.Poll,
	previousJudgment *db.Judgment,
	replaceMessage bool,
) error {
	return provider.GetResponder(input).RespondWithJudgmentUi(
		input,
		proposal,
		poll,
		previousJudgment,
		replaceMessage,
	)
}

//func RespondJudgmentSummary(
//	input provider.Input,
//	poll *db.Poll,
//	proposals []db.Proposal,
//	judgments []db.Judgment,
//	replaceMessage bool,
//) error {
//	return provider.GetResponder(input).RespondJudgmentSummary(
//		input,
//		poll,
//		proposals,
//		judgments,
//		replaceMessage,
//	)
//}

//func RespondBallotsInspection(
//	input provider.Input,
//	poll *db.Poll,
//	proposals []db.Proposal,
//	judgments []db.Judgment,
//) error {
//	return provider.GetResponder(input).RespondBallotsInspection(
//		input,
//		poll,
//		proposals,
//		judgments,
//	)
//}

func RespondServerError(input provider.Input, message string) error {
	return provider.GetResponder(input).RespondServerError(input, message)
}

func RespondUserError(input provider.Input, message string) error {
	return provider.GetResponder(input).RespondUserError(input, message)
}
