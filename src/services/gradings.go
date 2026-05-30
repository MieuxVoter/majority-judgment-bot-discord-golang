package services

import (
	"github.com/sarulabs/di/v2"
	"log"
	"log/slog"
	"main/src/container"
)

// Gradings provides the various gradings available to the polls.
// The gradings ABSOLUTELY MUST be unambiguously ordered.
type Gradings struct {
	logger *slog.Logger
}

// Get a grading.
func (service *Gradings) Get(key string) []string {
	// Right now these are hardcoded, but we could load them from config or something.
	switch key {
	case `👎👍`:
		return []string{"👎", "👍"}
	case `👎🤷👍`:
		return []string{"👎", "🤷", "👍"}
	case `🤮😐😀🤩`:
		return []string{"🤮", "😐", "😀", "🤩"}
	case `🤮😐😌😀🤩`:
		return []string{"🤮", "😐", "😌", "😀", "🤩"}
	}

	return service.Get(`🤮😐😌😀🤩`)
}

func GetGradings() *Gradings {
	return container.Get("gradings").(*Gradings)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "gradings",
		Build: func(ctn di.Container) (interface{}, error) {
			gradings := &Gradings{
				logger: ctn.Get("logger").(*slog.Logger),
			}
			return gradings, nil
		},
	})
	if err != nil {
		log.Fatalln("config failed to build", err)
	}
}
