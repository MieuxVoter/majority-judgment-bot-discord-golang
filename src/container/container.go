package container

// A simple Dependency Injection Container configuration.
// Make sure this module keeps no dependencies on other local modules.
//
// Dev Notes - Postmortem
// ----------------------
// DI is useful to collect collections of tagged services,
// but besides that, for now, it's not as useful as it is in other languages
// because of the way Go handles packages, which can serve as a makeshift DI.

import (
	"github.com/sarulabs/di"
	"log"
	"strings"
)

var builder *di.Builder
var container di.Container

// GetBuilder returns the container builder, to which we can add new services
func GetBuilder() *di.Builder {
	if builder == nil {
		var err error
		builder, err = di.NewBuilder()
		if err != nil {
			log.Fatalln(err)
		}
	}
	return builder
}

// Build the container ; done in main's init(), which is always ran last.
func Build() {
	container = GetBuilder().Build()
}

// Get a service, by name
func Get(name string) interface{} {
	return container.Get(name)
}

// GetCollection of services, by name prefix
func GetCollection(prefix string) []interface{} {
	collection := make([]interface{}, 0)
	for key := range container.Definitions() {
		if strings.HasPrefix(key, prefix+".") {
			collection = append(collection, container.Get(key))
		}
	}

	return collection
}
