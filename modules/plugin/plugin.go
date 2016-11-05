package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

// Plugin module handles the registration
// and management of plugins.
type Plugin struct {
	Plugins map[string]luaPlugin
	help    map[string]string
}

type luaPlugin struct {
	Name        string
	Description string
	Help        string
	Bind        map[string]*lua.LFunction
}

// NewPluginModule returns a new plugin module
func NewPluginModule() *Plugin {
	return &Plugin{
		Plugins: make(map[string]luaPlugin),
		help:    make(map[string]string),
	}
}

// Register will register the plugin module to
// the given lua state.
func (p *Plugin) Register(state *lua.LState) {
	state.PreloadModule("plugin", p.loader)
}

func (p *Plugin) loader(state *lua.LState) int {
	// Bind the functions to the module
	module := state.SetFuncs(state.NewTable(), map[string]lua.LGFunction{
		"add": p.Add,
	})

	// Push the module, exposing it to the Lua state
	state.Push(module)

	return 1
}

// Add will add a new plugin to the system
func (p *Plugin) Add(state *lua.LState) int {
	// Get the table
	table := state.ToTable(-1)

	// The registered plugin
	var plugin luaPlugin

	// Map the table
	if err := gluamapper.Map(table, &plugin); err != nil {
		log.Errorf("Could not parse plugin: %s", err)

		return 0
	}

	if _, ok := p.Plugins[plugin.Name]; ok {
		log.Errorf("Plugin with name '%s' already exist", plugin.Name)

		return 0
	}

	p.Plugins[plugin.Name] = plugin

	if plugin.Help != "" {
		p.help[plugin.Name] = plugin.Help
	}

	log.Infof("Registered plugin '%s'", plugin.Name)

	return 0
}

// Help will return the help string for the given function
// if it exists.
func (p *Plugin) Help(state *lua.LState) int {
	if val, ok := p.help[state.ToString(-1)]; ok && val != "" {
		state.Push(lua.LString(val))
	} else {
		state.Push(lua.LString("No help exists for that function"))
	}

	return 1
}
