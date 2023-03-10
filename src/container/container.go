package container

// A simple Dependency Injection Container configuration.
// Make sure this module keeps no dependencies on other local modules.
// This is required to keep our init() architecture.

import (
	"github.com/sarulabs/di"
	"log"
	"strings"
)

var builder *di.Builder
var container di.Container

func GetBuilder() *di.Builder {
	return builder
}

func Build() {
	container = GetBuilder().Build()
}

func Get(name string) interface{} {
	return container.Get(name)
}

func GetCollection(prefix string) []interface{} {
	collection := make([]interface{}, 0)
	for key, _ := range container.Definitions() {
		if strings.HasPrefix(key, prefix+".") {
			collection = append(collection, container.Get(key))
		}
	}

	return collection
}

func init() {
	//fmt.Println("init() container")
	var err error
	builder, err = di.NewBuilder()
	if err != nil {
		log.Fatalln(err)
	}
}
