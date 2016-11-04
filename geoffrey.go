package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/modules/geoffrey"
	"github.com/yuin/gopher-lua"
)

var (
	log = logrus.New()
)

func init() {
	// Output to stderr
	log.Out = os.Stderr

	// Set minimum log level
	log.Level = logrus.DebugLevel
}

func main() {
	state := lua.NewState(lua.Options{
		SkipOpenLibs: false,
	})

	// Close the Lua VM when we are done
	defer state.Close()

	// Add the geoffrey module
	geoffrey.Register(state)

	// Load config.lua
	if err := state.DoFile("config.lua"); err != nil {
		log.Fatalln(err)
	}

	// Load geoffrey.lua
	if err := state.DoFile("geoffrey.lua"); err != nil {
		log.Fatalln(err)
	}
}
