package database

type Poll struct {
	Id       uint64 `xorm:"pk autoincr"`
	AuthorId uint64 `xorm:"INDEX"`
	//Author   *disgord.User          `xorm:"-"`

	// The Subject of the poll should be somewhat short.
	Subject string
	//Subject string `xorm:"name"`

	//Gradation string `xorm:"-"`

	//DeadlineUnix timeutil.TimeStamp `xorm:"INDEX"`
	//CreatedUnix  timeutil.TimeStamp `xorm:"INDEX created"`
	//UpdatedUnix  timeutil.TimeStamp `xorm:"INDEX updated"`
	//ClosedUnix   timeutil.TimeStamp `xorm:"INDEX"`
}

func (poll *Poll) GetGradingSlice() []string {
	list := make([]string, 0, 5)

	// Placeholder until user customization somehow (poll.Grading?)
	// - 🤮😒😐🙂😀🤩
	// - 😫😒😐😌😀😍
	// - …
	list = append(list, "🤮")
	//list = append(list, "😒")
	list = append(list, "😐")
	list = append(list, "😌")
	list = append(list, "😀")
	list = append(list, "😍")

	return list
}

type Proposal struct {
	Id     uint64 `xorm:"pk autoincr"`
	Name   string
	PollId uint64 `xorm:"INDEX"`
}

type Judgment struct {
	JudgeSnowflake string `xorm:"INDEX(JX)"`
	ProposalId     uint64 `xorm:"INDEX(JX)"`
	PollId         uint64 `xorm:"INDEX(JX)"`
	Grade          uint8
}
