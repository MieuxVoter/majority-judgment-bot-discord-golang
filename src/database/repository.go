package database

import (
	"fmt"
	"github.com/andersfylling/disgord"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

func GetGuild(orm *xorm.Engine, snowflake disgord.Snowflake) (*Guild, error) {
	guild := &Guild{
		Snowflake: snowflake.String(),
	}
	has, err := orm.Get(guild)
	if !has {
		return nil, fmt.Errorf("no guild found for this snowflake")
	}
	if err != nil {
		return nil, err
	}

	return guild, nil
}

func CreateGuild(orm *xorm.Engine, snowflake disgord.Snowflake) (*Guild, error) {
	guild := &Guild{
		Snowflake: snowflake.String(),
		Quota:     42,
	}
	_, err := orm.InsertOne(guild)
	if err != nil {
		return nil, err
	}

	//logging.GetLogger().Warningln("New Guild !", guild.Snowflake)

	return guild, nil
}

func GetOrCreateGuild(orm *xorm.Engine, snowflake disgord.Snowflake) (guild *Guild, err error) {
	guild, err = GetGuild(orm, snowflake)
	if err != nil {
		guild, err = CreateGuild(orm, snowflake)
		if err != nil {
			return
		}
	}

	return
}

func GetPoll(orm *xorm.Engine, id uint64) (*Poll, error) {
	guild := &Poll{Id: id}
	has, err := orm.Get(guild)
	if !has {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return guild, nil
}

func GetPollProposals(e *xorm.Engine, poll *Poll) ([]Proposal, error) {
	var proposals []Proposal
	err := e.Where("poll_id = ?", poll.Id).Find(&proposals)
	if err != nil {
		return nil, err
	}

	return proposals, nil
}

func GetJudgmentsByJudgeOnPoll(e *xorm.Engine, judge *disgord.Member, poll *Poll) ([]Judgment, error) {
	if judge == nil {
		return nil, fmt.Errorf("no judge is defined")
	}

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

func CountGrades(e *xorm.Engine, poll *Poll, proposal *Proposal, gradeLevel uint8) (uint64, error) {
	rows := make([]int64, 0, 2)
	if err := e.Table("judgment").
		Select("COUNT(*) as amount").
		Where("`judgment`.`poll_id` = ?", poll.Id).
		And("`judgment`.`proposal_id` = ?", proposal.Id).
		And("`judgment`.`grade` = ?", gradeLevel).
		Find(&rows); err != nil {
		return 0, err
	}
	if 1 != len(rows) {
		return 0, fmt.Errorf("wrong shape in CountGrades")
	}

	return uint64(rows[0]), nil
}
