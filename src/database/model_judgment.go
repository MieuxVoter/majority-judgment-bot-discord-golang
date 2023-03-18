package database

// Judgment is the main citizen of this bot.
type Judgment struct {
	JudgeSnowflake string `xorm:"INDEX(JX) UNIQUE(JU)"`
	ProposalId     uint64 `xorm:"INDEX(JX) UNIQUE(JU)"`
	Grade          uint8  // 0 is "worst", most conservative, and usually default grade
}
