package main

import (
	"fmt"
	"os"

	"os/signal"
	"syscall"

	"io/ioutil"

	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/modules/geoffrey"
	"github.com/jriddick/geoffrey/modules/plugin"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

var (
	sigs = make(chan os.Signal, 1)
	bots *geoffrey.Geoffrey
	wg   sync.WaitGroup
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
		wg.Done()
	}()
}

type config struct {
	PluginFolder string
}

func main() {
	state := lua.NewState(lua.Options{
		SkipOpenLibs: false,
	})

	// Close the Lua VM when we are done
	defer state.Close()

	// Remove unsafe modules
	state.DoString("coroutine=nil;debug=nil;io=nil;os=nil")

	// Add the plugin module
	plugins := plugin.NewPluginModule()
	plugins.Register(state)

	// Add the geoffrey module
	bots := geoffrey.NewGeoffrey(plugins)
	bots.Register(state)

	// Load const.lua
	if err := state.DoFile("const.lua"); err != nil {
		log.WithError(err).Fatalln("Could not run const.lua")
	}

	// Load config.lua
	if err := state.DoFile("config.lua"); err != nil {
		log.WithError(err).Fatalln("Could not run config.lua")
	}

	// Map the configuration struct
	var cfg config
	if err := gluamapper.Map(state.GetGlobal("config").(*lua.LTable), &cfg); err != nil {
		log.WithError(err).Fatalln("Could not parse configuration struct")
	}

	// Read all plugins
	if files, err := ioutil.ReadDir(cfg.PluginFolder); err != nil {
		log.WithError(err).Fatalf("Could not open plugin directory '%s'", cfg.PluginFolder)
	} else {
		for _, file := range files {
			if !file.IsDir() {
				if err := state.DoFile(fmt.Sprintf("%s/%s", cfg.PluginFolder, file.Name())); err != nil {
					log.WithError(err).Errorf("Could not run file '%s'", file.Name())
				}
			}
		}
	}

	// Load geoffrey.lua
	if err := state.DoFile("geoffrey.lua"); err != nil {
		log.WithError(err).Fatalln("Could not run geoffrey.lua")
	}

	log.Infoln("Geoffrey is now running...")

	wg.Add(1)
	wg.Wait()
}
