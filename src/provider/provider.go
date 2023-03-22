package provider

import (
	"fmt"
	"main/src/container"
	"main/src/services"
)

// ResponderInterface should be implemented by our vendor output adapters
type ResponderInterface interface {
	Matches(input Input) bool
	RespondWithMessage(input Input, message string, ephemeral bool) error
	RespondWithMessageAndImage(input Input, message string, imageUrl string, ephemeral bool) error
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
