package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sarulabs/di"
	"main/src/configuration"
	"main/src/container"
	"xorm.io/xorm"
)

var Orm *xorm.Engine // deprecated

// Sync updates the database schema to mirror the models, if it can.
func Sync() error {
	return Engine().Sync(
		// We need to curate this list manually for now, sorry
		&Guild{},
		&Poll{},
		&Proposal{},
		&Judgment{},
	)
}

// Boot the ORN, configure it, store it and return it.
// Is idempotent, so it's safe to call this multiple times.
func Boot(config *configuration.Config) (*xorm.Engine, error) {
	if Orm != nil {
		return Orm, nil
	}

	databaseDriver := config.Get("DATABASE_DRIVER")
	databaseUrl := config.Get("DATABASE_URL")
	databaseCharset := config.Get("DATABASE_CHARSET")

	orm, err := xorm.NewEngine(databaseDriver, databaseUrl)
	if err != nil {
		return nil, err
	}

	orm.Charset(databaseCharset)

	if config.Get("APP_ENV") != "prod" {
		orm.ShowSQL(true)
	}

	Orm = orm

	return orm, nil
}

// Engine gets the currently booted ORM Engine, or nil
func Engine() *xorm.Engine {
	return container.Get("database.engine").(*xorm.Engine)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name:  "database.engine",
		Scope: di.App, // default
		Build: func(ctn di.Container) (interface{}, error) {
			engine, err := Boot(
				ctn.Get("config").(*configuration.Config),
			)
			if err != nil {
				return nil, err
			}
			return engine, err
		},
		Close: func(obj interface{}) error {
			return obj.(*xorm.Engine).Close()
		},
	})
	if err != nil {
		panic(err)
	}
}
