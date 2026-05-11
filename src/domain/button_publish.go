package domain

//
//import (
//	"github.com/sarulabs/di"
//	"log"
//	"main/src/container"
//	"main/src/provider"
//	"regexp"
//	"strconv"
//	"xorm.io/xorm"
//)
//
//var pollPublicationRegex = regexp.MustCompile("^button_publish:(?P<pollId>\\d+)$")
//
//type PublishButton struct {
//	orm *xorm.Engine
//}
//
//func (service PublishButton) Handle(input provider.Input) (handled bool, err error) {
//
//	handled = false
//	err = nil
//
//	buttonName, err := input.GetButtonName()
//	if err != nil {
//		return
//	}
//
//	matches := findNamedMatches(pollPublicationRegex, buttonName)
//	pollIdAsString, isMatchFound := matches["pollId"]
//
//	if !isMatchFound {
//		return
//	}
//
//	pollId, err := strconv.ParseUint(pollIdAsString, 10, 64)
//	if err != nil {
//		return false, err
//	}
//
//	handled, err = handlePollResult(service.orm, input, pollId, false)
//
//	return
//}
//
//func init() {
//	err := container.GetBuilder().Add(di.Def{
//		Name: "button.publish",
//		Build: func(ctn di.Container) (interface{}, error) {
//			cmd := &PublishButton{
//				orm: ctn.Get("database.engine").(*xorm.Engine),
//			}
//			return cmd, nil
//		},
//	})
//	if err != nil {
//		log.Fatalln("button.publish failed to build", err)
//	}
//}
