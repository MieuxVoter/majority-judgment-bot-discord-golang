package security

import (
	db "main/src/database"
	"xorm.io/xorm"
)

func CanGuildCreatePoll(_ *xorm.Engine, guild *db.Guild) (bool, error) {
	if guild == nil {
		return false, nil
	}

	if guild.Quota == 0 {
		return false, nil
	}

	return true, nil
}

func CanGuildParticipate(_ *xorm.Engine, _ *db.Guild) (bool, error) {
	// No bans for now
	return true, nil
}

func CanUserInspectBallots(_ *xorm.Engine, userVendorId string, poll *db.Poll) (bool, error) {
	if poll.Secrecy == "public" {
		return true, nil
	}
	return false, nil
}
