package database

import (
	"fmt"
	"time" // annoying IDE error, but benign
)

// Poll is our main (toplevel) model.
type Poll struct {
	Id              uint64 `xorm:"PK AUTOINCR"`
	GuildId         uint64 `xorm:"INDEX"`
	AuthorSnowflake string `xorm:"INDEX"`

	// Subject holds the purpose of the poll, to which proposals offer answers.
	// Eg: "The Next Meeting Date"
	// > What do you think of <Proposal> as <Subject> ?
	Subject string
	// Grading is a slice of unicode runes, as string.
	Grading string
	// Secrecy is either "public", "admin", or "secret".
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
