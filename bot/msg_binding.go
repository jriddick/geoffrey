package bot

import "github.com/yuin/gopher-lua"
import "github.com/jriddick/geoffrey/msg"

// RegisterMessage will register msg.Message struct to Lua
func RegisterMessage(state *lua.LState) {
	// Create the metatable
	meta := state.NewTypeMetatable("config")
	state.SetGlobal("message", meta)

	// Bind our functions
	state.SetField(meta, "__index", state.NewFunction(messageIndex))
}

// PushMessage will push the given msg.Message to the Lua stack
func PushMessage(msg *msg.Message, state *lua.LState) {
	// Create the userdata
	data := state.NewUserData()
	data.Value = msg

	// Set the metatable
	state.SetMetatable(data, state.GetTypeMetatable("message"))

	// Push the message
	state.Push(data)
}

func checkMessage(state *lua.LState) *msg.Message {
	data := state.CheckUserData(1)
	if v, ok := data.Value.(*msg.Message); ok {
		return v
	}
	state.ArgError(1, "message expected")
	return nil
}

func messageIndex(state *lua.LState) int {
	msg := checkMessage(state)
	key := state.CheckString(2)

	switch key {
	case "Command":
		state.Push(lua.LString(msg.Command))
	case "Params":
		params := state.NewTable()
		for key, param := range msg.Params {
			params.RawSetInt(key, lua.LString(param))
		}
		state.Push(params)
	case "Trailing":
		state.Push(lua.LString(msg.Trailing))
	default:
		state.Push(lua.LNil)
	}

	return 1
}
