package security

import (
	db "main/src/database"
	"xorm.io/xorm"
)

// CanGuildCreatePoll should be refactored as a service using DI
func CanGuildCreatePoll(_ *xorm.Engine, guild *db.Guild) (bool, error) {
	if guild == nil {
		return false, nil
	}

	if guild.Quota == 0 {
		return false, nil
	}

	return true, nil
}

// CanGuildParticipate is a security bouncer
func CanGuildParticipate(_ *xorm.Engine, _ *db.Guild) (bool, error) {
	// No bans for now
	return true, nil
}
