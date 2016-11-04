package geoffrey

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/bot"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

var (
	bots = make(map[string]*bot.Bot)
)

// Geoffrey module handles the creation, start, stop and
// cleanup of bots.
type Geoffrey struct {
}

// NewGeoffrey returns a new Geoffrey module
func NewGeoffrey() *Geoffrey {
	return &Geoffrey{}
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
	var config bot.Config

	// Map the configuration
	if err := gluamapper.Map(table, &config); err != nil {
		log.Fatalln(err)
	}

	// Create the bot
	bots[L.ToString(1)] = bot.NewBot(config)

	// Connect
	bots[L.ToString(1)].Connect()

	// Add logger handler
	bots[L.ToString(1)].OnMessage(func(bot *bot.Bot, channel *bot.Channel, user bot.User, msg string) {
		log.Debugln(msg)
	})

	// Run the handler
	bots[L.ToString(1)].Handler()

	return 0
}
