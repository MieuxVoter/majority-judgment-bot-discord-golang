package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

var Orm *xorm.Engine

func Sync() error {
	return Orm.Sync(
		// We need to curate this list manually for now, sorry
		&Poll{},
		&Proposal{},
		&Judgment{},
	)
}

func Boot(logLevel logrus.Level) (*xorm.Engine, error) {
	// todo: fetch these from env
	databaseDriver := "sqlite3"
	databaseName := "./mjbot.db"

	orm, err := xorm.NewEngine(databaseDriver, databaseName)
	if err != nil {
		return nil, err
	}

	if logLevel > logrus.InfoLevel {
		orm.ShowSQL(true)
	}

	Orm = orm

	return orm, nil
}
