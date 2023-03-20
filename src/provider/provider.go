package provider

import (
	"fmt"
	"main/src/container"
	"main/src/services"
)

type ResponderInterface interface {
	Matches(input Input) bool
	RespondWithMessage(input Input, message string, ephemeral bool) error
}

//type Responder struct {
//
//}

//func (r *Responder) Matches(input domain.Input, message string, ephemeral bool) error {
//
//}

//func (r *Responder) RespondWithMessage(input domain.Input, message string, ephemeral bool) error {
//
//}

func GetResponder(input Input) ResponderInterface {
	responders := container.GetCollection("responder")
	for _, genericResponder := range responders {
		responder := genericResponder.(ResponderInterface)
		if responder.Matches(input) {
			return responder
		}
	}

	services.GetLogger().Fatalln("no responder found")
	return nil
}

func RaiseInvalidProviderError(trace string) error {
	return fmt.Errorf("invalid input type for provider (%s)", trace)
}
