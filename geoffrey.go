package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/modules/geoffrey"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

func init() {
	// Output to stderr
	log.SetOutput(os.Stderr)
}

type config struct {
}

func main() {
	state := lua.NewState(lua.Options{
		SkipOpenLibs: false,
	})

	// Close the Lua VM when we are done
	defer state.Close()

	// Add the geoffrey module
	geoffrey.Register(state)

	// Load const.lua
	if err := state.DoFile("const.lua"); err != nil {
		log.Fatalln(err)
	}

	// Load config.lua
	if err := state.DoFile("config.lua"); err != nil {
		log.Fatalln(err)
	}

	// Map the configuration struct
	var cfg config
	if err := gluamapper.Map(state.GetGlobal("config").(*lua.LTable), &cfg); err != nil {
		log.Fatalln(err)
	}

	// Load geoffrey.lua
	if err := state.DoFile("geoffrey.lua"); err != nil {
		log.Fatalln(err)
	}
}
