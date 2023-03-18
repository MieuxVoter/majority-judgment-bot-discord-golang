package database

import (
	"fmt"
	"time" // annoying error, but benign
)

type Poll struct {
	Id              uint64 `xorm:"pk autoincr"`
	GuildId         uint64 `xorm:"INDEX"`
	AuthorSnowflake string `xorm:"INDEX"`

	Subject string
	// Grading is a slice of unicode runes
	Grading string
	// Secrecy is either "public", "admin", or "secret"
	Secrecy string

	CreatedUnix time.Time `xorm:"created"`
	UpdatedUnix time.Time `xorm:"updated"`
}

func (poll *Poll) getDefaultGrading() string {
	return "🤮😐😌😀🤩"
}

func (poll *Poll) GetGrading() string {
	if poll.Grading != "" {
		return poll.Grading
	}
	return poll.getDefaultGrading()
}

func (poll *Poll) GetGradingSlice() []string {
	list := make([]string, 0)
	for _, grade := range poll.GetGrading() {
		list = append(list, fmt.Sprintf("%c", grade))
	}

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

type Guild struct {
	Id        uint64 `xorm:"pk autoincr"`
	Snowflake string `xorm:"UNIQUE INDEX"`
	Quota     uint64
}
