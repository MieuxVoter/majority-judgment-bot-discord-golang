package old

//import (
//	"fmt"
//	//"github.com/andersfylling/disgord"
//	"github.com/sarulabs/di/v2"
//	"log"
//	"main/src/container"
//	db "main/src/database"
//	"main/src/provider"
//	"strconv"
//	"strings"
//	"xorm.io/xorm"
//)
//
//const RerunCommandSlug = "rerun"
//
//type RerunCommand struct {
//	orm *xorm.Engine
//}
//
//func (c RerunCommand) GetEmote() string {
//	return "♻"
//}
//
//func (c RerunCommand) GetName() string {
//	return RerunCommandSlug
//}
//
//func (c RerunCommand) GetDescription() string {
//	return "Rerun a fresh copy of a past poll"
//}
//
////func (c RerunCommand) DefineForDiscord() *disgord.ApplicationCommandOption {
////	return &disgord.ApplicationCommandOption{
////		Name:        c.GetName(),
////		Description: c.GetEmote() + " " + c.GetDescription(),
////		Type:        disgord.OptionTypeSubCommand,
////		Options: []*disgord.ApplicationCommandOption{
////			{
////				Name:        "poll",
////				Description: "The poll's numerical identifier, shown after the ⚖",
////				Type:        disgord.OptionTypeString,
////			},
////		},
////	}
////}
//
//func (c RerunCommand) Matches(command string) bool {
//	return command == c.GetName()
//}
//
//func (c RerunCommand) Handle(input provider.Input) (handled bool, err error) {
//	return true, handleRerunCommand(c.orm, input)
//}
//
//func handleRerunCommand(
//	orm *xorm.Engine,
//	input provider.Input,
//) error {
//
//	guildVendorId, _ := input.GetGuildVendorId()
//	guild, err := db.GetGuild(orm, guildVendorId)
//	if err != nil {
//		message := "This guild has no polls to rerun, and is not even registered yet.  " +
//			"Please create a poll first, using `/mj create …`."
//		return RespondUserError(input, message)
//	}
//
//	pollIdString, _ := input.GetOptionString("rerun", "poll", "")
//	if pollIdString == "" {
//		mostRecentPoll, errMrp := db.GetLastPollOfGuild(orm, guild)
//		if errMrp != nil {
//			message := "Fetching the most recent poll failed.  " +
//				"Please provide a poll identifier, or make sure you do have a poll to rerun."
//			return RespondUserError(input, message)
//		}
//		pollIdString = fmt.Sprint(mostRecentPoll.Id)
//	}
//	pollIdString = strings.Trim(pollIdString, "#")
//
//	pollId, errConv := strconv.Atoi(pollIdString)
//	if errConv != nil {
//		message := "🐉 The specified poll identifier is not a number.  " +
//			"Please use the _numerical_ identifier of the poll you want to _rerun_, " +
//			"like so `/mj rerun poll:42`."
//		return RespondUserError(input, message)
//	}
//	poll, errPoll := db.FindPoll(orm, uint64(pollId))
//	if errPoll != nil {
//		return errPoll
//	}
//	if poll == nil {
//		message := "The specified poll was not found.  " +
//			"Sorry.  _Here, have a banana instead : 🍌_"
//		return RespondUserError(input, message)
//	}
//	if poll.GuildId != guild.Id {
//		message := "The specified poll is in another castle. ⭐"
//		return RespondUserError(input, message)
//	}
//
//	proposals, errProp := db.GetPollProposals(orm, poll)
//	if errProp != nil {
//		return errProp
//	}
//
//	subject := poll.Subject
//	proposalsNames := make([]string, 0)
//	for _, proposal := range proposals {
//		proposalsNames = append(proposalsNames, proposal.Name)
//	}
//	grading := poll.Grading
//	secrecy := poll.Secrecy
//
//	err = doCreatePoll(orm, input, subject, proposalsNames, grading, secrecy)
//
//	return err
//}
//
//func init() {
//	err := container.GetBuilder().Add(di.Def{
//		Name: "command." + RerunCommandSlug,
//		Build: func(ctn di.Container) (interface{}, error) {
//			cmd := &RerunCommand{
//				orm: ctn.Get("database.engine").(*xorm.Engine),
//			}
//			return cmd, nil
//		},
//	})
//	if err != nil {
//		log.Fatalf("command.%s failed to build : %s\n", ExplainCommandSlug, err)
//	}
//}
