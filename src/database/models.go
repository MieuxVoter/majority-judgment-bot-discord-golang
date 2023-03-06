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

	// No idea how xorm works -- help!
	// Judgments         []*PollJudgment    `xorm:"-"`
	// Judgments         JudgmentList   `xorm:"-"`
}

type Proposal struct {
	Id     uint64 `xorm:"pk autoincr"`
	Name   string
	PollId uint64 `xorm:"INDEX"`
}
