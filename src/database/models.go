package database

import "time"

type Poll struct {
	Id       uint64 `xorm:"pk autoincr"`
	AuthorId uint64 `xorm:"INDEX"`
	//Author   *disgord.User          `xorm:"-"`

	// The Subject of the poll should be somewhat short.
	Subject string
	//Subject string `xorm:"name"`

	//Gradation string `xorm:"-"`

	CreatedUnix time.Time `xorm:"created"`
	UpdatedUnix time.Time `xorm:"updated"`
	//DeadlineUnix timeutil.TimeStamp `xorm:"INDEX"`
	//ClosedUnix   timeutil.TimeStamp `xorm:"INDEX"`
}

func (poll *Poll) GetGradingSlice() []string {
	list := make([]string, 0, 5)

	// Placeholder shim until user customization somehow (poll.Grading?)
	// Careful: for now only 5 grades max are supported.  (amount of buttons per action row)
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

func (poll *Poll) GetGradeIcon(gradeLevel uint8) string {
	icons := poll.GetGradingSlice()
	gradeLevelInt := int(gradeLevel)
	if len(icons) <= gradeLevelInt {
		return "🥚" // easter
	}

	return icons[gradeLevelInt]
}

type Proposal struct {
	Id     uint64 `xorm:"pk autoincr"`
	Name   string
	PollId uint64 `xorm:"INDEX"`
}

type Judgment struct {
	JudgeSnowflake string `xorm:"INDEX(JX) UNIQUE(JU)"`
	ProposalId     uint64 `xorm:"INDEX(JX) UNIQUE(JU)"`
	PollId         uint64 `xorm:"INDEX(JX) UNIQUE(JU)"`
	Grade          uint8
}
