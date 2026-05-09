package domain

import (
	"main/src/provider"
	"regexp"
)

type Button interface {
	// GetRegex is not strictly required anymore, at least not for Discord, as the pattern is enough to identify the button,
	// and we could grab the poll id from the event data.  It's fine like this, though.
	GetRegex() *regexp.Regexp
	GetPattern() string
	Handle(input provider.ButtonInput) (bool, error)
}

// TODO: move this to utils package or file?
func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}

	return results
}
