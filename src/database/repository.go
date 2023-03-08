package database

import (
	"github.com/andersfylling/disgord"
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

func GetJudgmentsByJudgeOnPoll(e *xorm.Engine, judge *disgord.Member, poll *Poll) ([]Judgment, error) {
	var judgments []Judgment
	err := e.
		Where("judge_snowflake = ?", judge.UserID.String()).
		Where("poll_id = ?", poll.Id).
		OrderBy("proposal_id", "ASC").
		Find(&judgments)
	if err != nil {
		return nil, err
	}

	return judgments, nil
}
