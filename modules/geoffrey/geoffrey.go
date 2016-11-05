package geoffrey

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
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

// Config is the configuration struct for Geoffrey
type Config struct {
	Connection struct {
		Host               string
		Port               int
		Secure             bool
		InsecureSkipVerify bool
		Timeout            int
	}
	Limits struct {
		Reconnect int
		Timeout   int
	}
	Authentication struct {
		Username string
		Password string
	}
	Registration struct {
		Nick string
		User string
		Name string
	}
	Channels []string
	Plugins  []string
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
func (g *Geoffrey) Add(state *lua.LState) int {
	// Build the logger
	logger := log.WithField("file", state.Where(1))

	// Check so we got the required number of parameters
	if state.GetTop() != 2 {
		logger.Errorf("Geoffrey:Add takes two parameters, we got %d", state.GetTop())
		return 0
	}

	// Get the table
	if state.Get(2).Type() != lua.LTTable {
		logger.Errorf("Geoffrey:Add takes a table as a second parameter, we got '%s'", state.Get(2).Type().String())
		return 0
	}
	table := state.Get(2).(*lua.LTable)

	// Get the name
	if state.Get(1).Type() != lua.LTString {
		logger.Errorf("Geoffrey:Add takes a string as the first parameter, we got '%s'", state.Get(1).Type().String())
		return 0
	}
	name := lua.LVAsString(state.Get(1))

	// The bot configuration
	var config Config

	// Map the configuration
	if err := gluamapper.Map(table, &config); err != nil {
		log.Fatalln(err)
	}

	// Create the bot
	bots[name] = bot.NewBot(bot.Config{
		Hostname:           config.Connection.Host,
		Port:               config.Connection.Port,
		Secure:             config.Connection.Secure,
		InsecureSkipVerify: config.Connection.InsecureSkipVerify,
		Timeout:            config.Connection.Timeout,
		TimeoutLimit:       config.Limits.Timeout,
		ReconnectLimit:     config.Limits.Reconnect,
		Nick:               config.Registration.Nick,
		User:               config.Registration.User,
		Name:               config.Registration.Name,
		Channels:           config.Channels,
	}, state.NewThread())

	// Get all loaded plugins
	for _, pluginName := range config.Plugins {
		if plugin, ok := g.plugins.Plugins[pluginName]; !ok {
			logger.Errorf("Tried to use non-existent plugin '%s'", pluginName)
			continue
		} else {
			for event, handler := range plugin.Bind {
				switch event {
				case "OnMessage":
					event = irc.Message
				case "OnPing":
					event = irc.Ping
				case "OnWelcome":
					event = irc.Welcome
				case "OnJoin":
					event = irc.Join
				case "OnPart":
					event = irc.Part
				case "OnNameList":
					event = irc.Namreply
				case "OnNotice":
					event = irc.Notice
				}

				bots[name].AddLuaHandler(event, handler)
			}
		}
	}

	// Connect
	if err := bots[name].Connect(); err != nil {
		logger.WithError(err).Errorln("We could not connect to the server")
	}

	// Run the handler
	go bots[name].Handler()

	return 0
}
