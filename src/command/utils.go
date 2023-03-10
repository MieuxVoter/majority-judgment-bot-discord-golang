package command

import (
	"fmt"
	"github.com/andersfylling/disgord"
)

func getSubcommandOptions(
	options []*disgord.ApplicationCommandDataOption,
	name string) ([]*disgord.ApplicationCommandDataOption, error) {

	for _, option := range options {
		if option.Name == name {
			return option.Options, nil
		}
	}

	return nil, fmt.Errorf("command subquery not found")
}

func getOptionStringByName(
	options []*disgord.ApplicationCommandDataOption,
	name string,
	defaultValue string) string {

	for _, option := range options {
		if option.Name == name {
			value := fmt.Sprintf("%s", option.Value)
			if value == "" {
				value = defaultValue
			}

			return value
		}
	}

	return defaultValue
}
