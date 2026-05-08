package container

import (
	"github.com/sarulabs/di"
	"log"
	"strings"
)

// A simple Dependency Injection Container configuration.
// Make sure this module keeps no dependencies on other local modules.
//
// Add your services to the builder in init() functions.
// I know it's usually considered bad practice to use init() in Go,
// but I'm willing to make an exception for the DIC.
// Do not add anything else to init() but your service definitions, por favor.
//
// In addition, remember that Go has aggressive tree shaking,
// so your init() might not run if your file is never imported anywhere.
//
// Dev Notes - Postmortem
// ----------------------
// DI is actually useful to collect collections of tagged services,
// but besides that, for now, it's not as useful as it is in other languages
// because of the way Go handles packages, which can serve as a makeshift DI.

var builder *di.Builder
var container di.Container

// Build the container ; called in main's init(), which always runs last.
// We need to add all our services to the builder *before* we build.
func Build() {
	container = GetBuilder().Build()
}

// GetBuilder returns the container builder, to which we can add new services.
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

// Get a service, by name.
func Get(name string) interface{} {
	return container.Get(name)
}

// GetCollection of services, by name prefix.
func GetCollection(prefix string) []interface{} {
	collection := make([]interface{}, 0)
	for key := range container.Definitions() {
		if strings.HasPrefix(key, prefix+".") {
			collection = append(collection, container.Get(key))
		}
	}

	return collection
}
