package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"
	"xorm.io/xorm"
)

var Orm *xorm.Engine

// Sync updates the database schema to mirror the models, if it can.
func Sync() error {
	return Orm.Sync(
		// We need to curate this list manually for now, sorry
		&Guild{},
		&Poll{},
		&Proposal{},
		&Judgment{},
	)
}

// Boot the ORN, configure it, store it and return it.
// Is idempotent, so it's safe to call this multiple times.
func Boot(logLevel logrus.Level) (*xorm.Engine, error) {
	if Orm != nil {
		return Orm, nil
	}

	databaseDriver := os.Getenv("DATABASE_DRIVER")
	databaseUrl := os.Getenv("DATABASE_URL")
	databaseCharset := os.Getenv("DATABASE_CHARSET")

	orm, err := xorm.NewEngine(databaseDriver, databaseUrl)
	if err != nil {
		return nil, err
	}

	orm.Charset(databaseCharset)

	if logLevel > logrus.InfoLevel {
		orm.ShowSQL(true)
	}

	Orm = orm

	return orm, nil
}

// Get the currently booted ORM Engine, or nil
func Engine() *xorm.Engine {
	return Orm
}
