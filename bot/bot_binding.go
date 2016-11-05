package bot

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/helper"
	"github.com/yuin/gopher-lua"
)

// RegisterBot will register the Bot struct to Lua
func RegisterBot(state *lua.LState) {
	// Create the Metatable
	meta := state.NewTypeMetatable("bot")
	state.SetGlobal("bot", meta)

	// Bind our functions
	state.SetField(meta, "__index", state.NewFunction(botIndex))
}

// PushBot will push an existing *Bot onto the Lua stack
func PushBot(bot *Bot, state *lua.LState) {
	// Create the bot user data
	data := state.NewUserData()
	data.Value = bot

	// Set the Metatable
	state.SetMetatable(data, state.GetTypeMetatable("bot"))

	// Push the bot
	state.Push(data)
}

func checkBot(state *lua.LState) *Bot {
	// Get the logger
	logger := log.WithField("file", state.Where(1))

	// Try to get the userdata
	data := helper.GetUserData(1, state)
	if data != nil {
		// Check if userdata is the right type
		if bot, ok := data.Value.(*Bot); ok {
			return bot
		}

		logger.Errorf("Expected userdata type of of type Bot, we got '%s'", reflect.TypeOf(data.Value).Name())
	}
	return nil
}

func botSend(state *lua.LState) int {
	bot := checkBot(state)
	rcv := helper.GetString(2, state)
	msg := helper.GetString(3, state)

	if msg != nil && rcv != nil {
		bot.Send(*rcv, *msg)
	}

	return 0
}

func botJoin(state *lua.LState) int {
	bot := checkBot(state)
	channel := helper.GetString(2, state)

	if channel != nil {
		bot.Join(*channel)
	}

	return 0
}

func botPing(state *lua.LState) int {
	bot := checkBot(state)
	msg := helper.GetString(2, state)

	if msg != nil {
		bot.Ping(*msg)
	}

	return 0
}

func botPong(state *lua.LState) int {
	bot := checkBot(state)
	msg := helper.GetString(2, state)

	if msg != nil {
		bot.Pong(*msg)
	}

	return 0
}

func botNick(state *lua.LState) int {
	bot := checkBot(state)
	nick := helper.GetString(2, state)

	if nick != nil {
		bot.Nick(*nick)
	}

	return 0
}

func botUser(state *lua.LState) int {
	bot := checkBot(state)
	user := helper.GetString(2, state)
	name := helper.GetString(3, state)

	if user != nil && name != nil {
		bot.User(*user, *name)
	}

	return 0
}

func botIndex(state *lua.LState) int {
	bot := checkBot(state)

	if bot == nil {
		return 0
	}

	key := helper.GetString(2, state)

	if key == nil {
		return 0
	}

	switch *key {
	case "send":
		state.Push(state.NewFunction(botSend))
	case "join":
		state.Push(state.NewFunction(botJoin))
	case "ping":
		state.Push(state.NewFunction(botPing))
	case "pong":
		state.Push(state.NewFunction(botPong))
	case "nick":
		state.Push(state.NewFunction(botNick))
	case "user":
		state.Push(state.NewFunction(botUser))
	case "config":
		PushConfig(&bot.config, state)
	default:
		state.Push(lua.LNil)
	}

	return 1
}
