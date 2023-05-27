package database

// Feedback holds a feedback sent by a user
type Feedback struct {
	Id             uint64 `xorm:"PK AUTOINCR"`
	GuildId        uint64
	Content        string
	AuthorName     string
	AuthorVendorId string
}
