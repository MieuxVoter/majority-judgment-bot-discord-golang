package database

import (
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

func GetPollProposals(e *xorm.Engine, poll *Poll) ([]Proposal, error) {
	var proposals []Proposal
	err := e.Where("poll_id = ?", poll.Id).Find(&proposals)
	if err != nil {
		return nil, err
	}

	return proposals, nil
}
