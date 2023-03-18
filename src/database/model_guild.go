package database

// Guild represents a community that invoked the bot.
type Guild struct {
	Id        uint64 `xorm:"pk autoincr"`
	Snowflake string `xorm:"UNIQUE INDEX"`
	Quota     uint64 // Remaining amount of Polls this guild is allowed to create.
}
