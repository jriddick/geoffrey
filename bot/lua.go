package bot

import "github.com/yuin/gopher-lua"

// RegisterBot will register the Bot struct to Lua
func RegisterBot(state *lua.LState) {
	// Create the Metatable
	meta := state.NewTypeMetatable("bot")
	state.SetGlobal("bot", meta)

	// Bind our functions
	state.SetField(meta, "__index", state.NewFunction(botIndex))
}

// RegisterConfig will register the Config struct to Lua
func RegisterConfig(state *lua.LState) {
	// Create the Metatable
	meta := state.NewTypeMetatable("config")
	state.SetGlobal("config", meta)

	// Bind our functions
	state.SetField(meta, "__index", state.NewFunction(configIndex))
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

// PushConfig will push the given config to the Lua stack
func PushConfig(config *Config, state *lua.LState) {
	// Create the config user data
	data := state.NewUserData()
	data.Value = config

	// Set the Metatable
	state.SetMetatable(data, state.GetTypeMetatable("config"))

	// Push the config
	state.Push(data)
}

// checkBot will check wether the first argument is a *Bot
// or not.
func checkBot(state *lua.LState) *Bot {
	data := state.CheckUserData(1)
	if v, ok := data.Value.(*Bot); ok {
		return v
	}
	state.ArgError(1, "bot expected")
	return nil
}

// checkConfig will check wether the first argument is of type *Config
// and return it.
func checkConfig(state *lua.LState) *Config {
	data := state.CheckUserData(1)
	if v, ok := data.Value.(*Config); ok {
		return v
	}
	state.ArgError(1, "config expected")
	return nil
}

func botSend(state *lua.LState) int {
	bot := checkBot(state)
	rcv := state.CheckString(2)
	msg := state.CheckString(3)

	bot.Send(rcv, msg)
	return 0
}

func botJoin(state *lua.LState) int {
	bot := checkBot(state)
	channel := state.CheckString(2)

	bot.Join(channel)
	return 0
}

func botIndex(state *lua.LState) int {
	bot := checkBot(state)
	key := state.CheckString(2)

	switch key {
	case "send":
		state.Push(state.NewFunction(botSend))
	case "join":
		state.Push(state.NewFunction(botJoin))
	case "config":
		PushConfig(&bot.config, state)
	default:
		state.Push(lua.LNil)
	}

	return 1
}

func configIndex(state *lua.LState) int {
	config := checkConfig(state)
	key := state.CheckString(2)

	switch key {
	case "Hostname":
		state.Push(lua.LString(config.Hostname))
	case "Port":
		state.Push(lua.LNumber(config.Port))
	case "Secure":
		state.Push(lua.LBool(config.Secure))
	case "InsecureSkipVerify":
		state.Push(lua.LBool(config.InsecureSkipVerify))
	case "Nick":
		state.Push(lua.LString(config.Nick))
	case "User":
		state.Push(lua.LString(config.User))
	case "Name":
		state.Push(lua.LString(config.Name))
	case "Channels":
		channels := state.NewTable()
		for key, value := range config.Channels {
			channels.RawSetInt(key, lua.LString(value))
		}
		state.Push(channels)
	case "Timeout":
		state.Push(lua.LNumber(config.Timeout))
	case "TimeoutLimit":
		state.Push(lua.LNumber(config.TimeoutLimit))
	case "ReconnectLimit":
		state.Push(lua.LNumber(config.ReconnectLimit))
	default:
		state.Push(lua.LNil)
	}

	return 1
}
