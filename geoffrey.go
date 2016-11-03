package main

import (
	"log"

	"github.com/jriddick/geoffrey/modules/geoffrey"
	"github.com/yuin/gopher-lua"
)

func main() {
	state := lua.NewState(lua.Options{
		SkipOpenLibs: false,
	})

	// Close the Lua VM when we are done
	defer state.Close()

	// Add the geoffrey module
	state.PreloadModule("geoffrey", geoffrey.Load)

	// Load config.lua
	if err := state.DoFile("config.lua"); err != nil {
		log.Fatalln(err)
	}
}
