package commands

import (
	"github.com/sarulabs/di/v2"
	"log"
	"main/src/container"
)

// MjCommand is our main (and only) root slash command for now.
// It does not do anything by itself, and instead relies on subcommands.
type MjCommand struct{}

func (c MjCommand) GetName() string {
	return "mj"
}

func (c MjCommand) GetDescription() string {
	return "Manage Majority Judgment polls"
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "command.mj",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := MjCommand{}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalf("service command.mj failed to build: %s\n", err)
	}
}
