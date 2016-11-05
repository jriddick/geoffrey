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
	// Build the logger
	logger := log.WithField("file", state.Where(1))

	// Check so we get the required number of parameters
	if state.GetTop() != 1 {
		logger.Errorf("Plugin:Add takes one parameter, we got '%d'", state.GetTop())
		return 0
	}

	// Get the table
	if state.Get(1).Type() != lua.LTTable {
		logger.Errorf("Geoffrey:Add takes table as a parameter, we got '%s'", state.Get(1).Type().String())
		return 0
	}
	table := state.Get(1).(*lua.LTable)

	// The registered plugin
	var plugin luaPlugin

	// Map the table
	if err := gluamapper.Map(table, &plugin); err != nil {
		logger.WithError(err).Errorln("Could not parse the plugin")
		return 0
	}

	if _, ok := p.Plugins[plugin.Name]; ok {
		logger.Errorf("Plugin with name '%s' already exist", plugin.Name)
		return 0
	}

	p.Plugins[plugin.Name] = plugin

	if plugin.Help != "" {
		p.help[plugin.Name] = plugin.Help
	}

	logger.Infof("Registered plugin '%s'", plugin.Name)

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
