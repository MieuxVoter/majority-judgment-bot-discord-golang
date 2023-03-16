package database

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

func GetGuild(orm *xorm.Engine, snowflake string) (*Guild, error) {
	guild := &Guild{
		Snowflake: snowflake,
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

func CreateGuild(orm *xorm.Engine, snowflake string) (*Guild, error) {
	guild := &Guild{
		Snowflake: snowflake,
		Quota:     111,
	}
	_, err := orm.InsertOne(guild)
	if err != nil {
		return nil, err
	}

	//logging.GetLogger().Warningln("New Guild !", guild.Snowflake)

	return guild, nil
}

func GetOrCreateGuild(orm *xorm.Engine, snowflake string) (guild *Guild, err error) {
	guild, err = GetGuild(orm, snowflake)
	if err != nil {
		guild, err = CreateGuild(orm, snowflake)
		if err != nil {
			return
		}
	}

	return
}

// FindPoll returns nil if the poll was not found
func FindPoll(orm *xorm.Engine, id uint64) (*Poll, error) {
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

// CountPolls counts all the polls this bot is responsible for
func CountPolls(orm *xorm.Engine) (int64, error) {
	poll := &Poll{}
	count, err := orm.Count(poll)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountGuildPolls counts all the polls this bot is responsible for
func CountGuildPolls(orm *xorm.Engine, guild *Guild) (int64, error) {
	poll := &Poll{}
	count, err := orm.Where("guild_id = ?", guild.Id).Count(poll)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetLastPollOfGuild returns the most recent poll of the specified guild, or fails.
func GetLastPollOfGuild(orm *xorm.Engine, guild *Guild) (*Poll, error) {
	poll := &Poll{}
	has, err := orm.
		Where("guild_id = ?", guild.Id).
		OrderBy("created_unix", "DESC").
		Limit(1).
		Get(poll)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("no last poll found")
	}

	return poll, nil
}

func GetPollProposals(e *xorm.Engine, poll *Poll) ([]Proposal, error) {
	var proposals []Proposal
	err := e.Where("poll_id = ?", poll.Id).Find(&proposals)
	if err != nil {
		return nil, err
	}

	return proposals, nil
}

func GetJudgmentsByJudgeOnPoll(orm *xorm.Engine, judge string, poll *Poll) ([]Judgment, error) {
	if judge == "" {
		return nil, fmt.Errorf("no judge is defined")
	}

	var judgments []Judgment
	err := orm.
		Where("judge_snowflake = ?", judge).
		Where("poll_id = ?", poll.Id).
		OrderBy("proposal_id", "ASC").
		Find(&judgments)
	if err != nil {
		return nil, err
	}

	return judgments, nil
}

func CountGrades(orm *xorm.Engine, poll *Poll, proposal *Proposal, gradeLevel uint8) (uint64, error) {
	rows := make([]int64, 0, 2)
	if err := orm.Table("judgment").
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

func CollectAllJudgmentsOnPoll(e *xorm.Engine, poll *Poll, proposals []Proposal) ([]Judgment, error) {
	var proposalsIds = make([]uint64, 0)
	for _, proposal := range proposals {
		proposalsIds = append(proposalsIds, proposal.Id)
	}

	var judgments []Judgment
	err := e.
		In("proposal_id", proposalsIds).
		OrderBy("judge_snowflake", "ASC").
		OrderBy("proposal_id", "ASC").
		Find(&judgments)
	if err != nil {
		return nil, err
	}

	return judgments, nil
}
