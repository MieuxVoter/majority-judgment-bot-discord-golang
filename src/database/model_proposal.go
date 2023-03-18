package database

// Proposal represents a single Proposal (aka. Candidate) of a Poll.
type Proposal struct {
	Id     uint64 `xorm:"PK AUTOINCR"`
	Name   string
	PollId uint64 `xorm:"INDEX"`
}
