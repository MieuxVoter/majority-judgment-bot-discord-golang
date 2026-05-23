package old

//
//import (
//	"github.com/sarulabs/di/v2"
//	"log"
//	"main/src/container"
//	db "main/src/database"
//	"main/src/provider"
//	"regexp"
//	"strconv"
//	"xorm.io/xorm"
//)
//
//var pollInspectionRegex = regexp.MustCompile("^button_inspect:(?P<pollId>\\d+)$")
//
//type InspectButton struct {
//	orm *xorm.Engine
//}
//
//func (service InspectButton) Handle(input provider.Input) (handled bool, err error) {
//
//	handled = false
//	err = nil
//
//	buttonName, _ := input.GetButtonName()
//	matches := findNamedMatches(pollInspectionRegex, buttonName)
//	pollIdAsString, isMatchFound := matches["pollId"]
//
//	if !isMatchFound {
//		return
//	}
//
//	pollId, errParse := strconv.ParseUint(pollIdAsString, 10, 64)
//	if errParse != nil {
//		return false, errParse
//	}
//
//	poll := &db.Poll{Id: pollId}
//	foundPoll, errOrm := service.orm.Get(poll)
//	if errOrm != nil {
//		err = RespondServerError(input, "Ooops.  This poll was probably deleted.")
//		return
//	}
//	if !foundPoll {
//		err = RespondServerError(input, "Oh noes!  This poll was probably deleted.")
//		return
//	}
//
//	handled = true
//
//	proposals, err := db.GetPollProposals(service.orm, poll)
//	if err != nil {
//		return false, err
//	}
//	if len(proposals) == 0 {
//		err = RespondServerError(input, "Wait a minute…  This poll has no proposals !?")
//		return
//	}
//
//	allJudgments, errCollect := db.CollectAllJudgmentsOnPoll(service.orm, poll, proposals)
//	if errCollect != nil {
//		err = RespondServerError(input, "I somehow cannot collect judgments for this poll.  🔔")
//		return
//	}
//
//	err = RespondBallotsInspection(input, poll, proposals, allJudgments)
//	return
//}
//
//func init() {
//	err := container.GetBuilder().Add(di.Def{
//		Name: "button.inspect",
//		Build: func(ctn di.Container) (interface{}, error) {
//			cmd := &InspectButton{
//				orm: ctn.Get("database.engine").(*xorm.Engine),
//			}
//			return cmd, nil
//		},
//	})
//	if err != nil {
//		log.Fatalln("button.inspect failed to build", err)
//	}
//}
