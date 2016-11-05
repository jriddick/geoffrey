package geoffrey

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/modules/plugin"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

var (
	bots = make(map[string]*bot.Bot)
)

// Geoffrey module handles the creation, start, stop and
// cleanup of bots.
type Geoffrey struct {
	plugins *plugin.Plugin
}

type geoffreyConfig struct {
	Config  bot.Config
	Plugins []string
}

// NewGeoffrey returns a new Geoffrey module
func NewGeoffrey(plugins *plugin.Plugin) *Geoffrey {
	return &Geoffrey{
		plugins: plugins,
	}
}

// Register will register the module to the lua state
func (g *Geoffrey) Register(L *lua.LState) {
	L.RegisterModule("geoffrey", map[string]lua.LGFunction{
		"add": g.Add,
	})
}

// Shutdown will stop all running bots
func (g *Geoffrey) Shutdown() {
	for _, bot := range bots {
		bot.Close()
	}
}

// Add will create and start a new bot
func (g *Geoffrey) Add(L *lua.LState) int {
	// Get the table
	table := L.ToTable(-1)

	// The bot configuration
	var config geoffreyConfig

	// Map the configuration
	if err := gluamapper.Map(table, &config); err != nil {
		log.Fatalln(err)
	}

	// Create the bot
	bots[L.ToString(1)] = bot.NewBot(config.Config, L.NewThread())

	// Get all loaded plugins
	for _, name := range config.Plugins {
		if plugin, ok := g.plugins.Plugins[name]; !ok {
			log.Errorf("Tried to use non-existent plugin '%s'", name)
			continue
		} else {
			if plugin.Bind.OnMessage != nil {
				// Bind the plugin handler
				bots[L.ToString(1)].OnMessageLua(plugin.Bind.OnMessage)
			}
		}
	}

	// Connect
	bots[L.ToString(1)].Connect()

	// Run the handler
	go bots[L.ToString(1)].Handler()

	return 0
}
