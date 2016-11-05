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
	key := state.CheckString(2)

	switch key {
	case "send":
		state.Push(state.NewFunction(botSend))
	case "join":
		state.Push(state.NewFunction(botJoin))
	default:
		state.Push(lua.LNil)
	}

	return 1
}
