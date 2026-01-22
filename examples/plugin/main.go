package main

import (
	"github.com/reglet-dev/reglet-sdk/go/application/plugin"
)

func init() {
	// Register the plugin with the SDK lifecycle
	plugin.Register(&ExamplePlugin{})
}

func main() {
	// main is not called in -buildmode=c-shared
}
