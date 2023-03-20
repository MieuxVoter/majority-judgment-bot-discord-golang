package domain

import (
	"main/src/provider"
	"regexp"
)

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

type Button interface {
	Handle(input provider.Input) (bool, error)
}

//func HandleButtonDeliberate(
//	ctx context.Context,
//	s disgord.Session,
//	h *disgord.InteractionCreate,
//) (handled bool, err error) {
//
//}
