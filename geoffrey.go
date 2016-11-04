package main

import (
	"os"

	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/modules/geoffrey"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

var (
	sigs = make(chan os.Signal, 1)
	bots *geoffrey.Geoffrey
)

func init() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Set the log level to debug
	log.SetLevel(log.DebugLevel)

	// Capture signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		// Wait until we get a signal
		<-sigs

		// Shutdown the bot
		bots.Shutdown()

		// Exit the program
		os.Exit(1)
	}()
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
	bots := geoffrey.NewGeoffrey()
	bots.Register(state)

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

	log.Println("geoffrey is now running")

	// Block until program exists (Ctrl-C)
	for {
	}
}
