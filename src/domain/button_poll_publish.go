package domain

import (
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	"main/src/provider"
	"regexp"
	"strconv"
	"xorm.io/xorm"
)

var buttonPollPublishRegex = regexp.MustCompile("^/button/poll/(?P<pollId>\\d+)/publish$")
var buttonPollPublishPattern = "/button/poll/{pollId}/publish"

type PollPublishButton struct {
	orm *xorm.Engine
}

func (b PollPublishButton) GetRegex() *regexp.Regexp {
	return buttonPollPublishRegex
}

func (b PollPublishButton) GetPattern() string {
	return buttonPollPublishPattern
}

func (b PollPublishButton) Handle(input provider.ButtonInput) (handled bool, err error) {

	handled = false
	err = nil

	buttonName, err := input.GetButtonName()
	if err != nil {
		return
	}

	matches := findNamedMatches(buttonPollPublishRegex, buttonName)
	pollIdAsString, isMatchFound := matches["pollId"]

	if !isMatchFound {
		return
	}

	pollId, err := strconv.ParseUint(pollIdAsString, 10, 64)
	if err != nil {
		return
	}

	handled, err = handlePollResult(b.orm, input, pollId, false)

	return
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "button.poll.publish",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &PollPublishButton{
				orm: ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalln("button.poll.publish failed to build", err)
	}
}
