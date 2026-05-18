package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
	"main/src/services"
	"xorm.io/xorm"
)

// GetEngine gets the currently booted ORM Engine, or nil
func GetEngine() *xorm.Engine {
	return container.Get("database.engine").(*xorm.Engine)
}

// Sync updates the database schema to mirror the Go models, if it can.
func Sync() error {
	return GetEngine().Sync(
		// We need to curate this list manually for now, sorry.  (-_-)
		// We could use init() shenanigans like Gitea does to fix this.
		&Guild{},
		&Poll{},
		&Proposal{},
		&Judgment{},
		&Feedback{},
	)
}

// bootEngine boots an ORM according to passed configuration, and returns it.
func bootEngine(config *services.Config) (*xorm.Engine, error) {
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

	return orm, nil
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name:  "database.engine",
		Scope: di.App, // default
		Build: func(ctn di.Container) (interface{}, error) {
			engine, err := bootEngine(
				ctn.Get("config").(*services.Config),
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
		log.Fatal(err)
	}
}
