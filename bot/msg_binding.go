package bot

import "github.com/yuin/gopher-lua"
import "github.com/jriddick/geoffrey/msg"

// RegisterMessage will register msg.Message struct to Lua
func RegisterMessage(state *lua.LState) {
	// Create the metatable
	meta := state.NewTypeMetatable("message")
	state.SetGlobal("message", meta)

	// Bind the message index function
	state.SetField(meta, "__index", state.NewFunction(messageIndex))

	// Create the metatable for prefix struct
	prefix := state.NewTypeMetatable("prefix")
	state.SetGlobal("prefix", prefix)

	// Bind the prefix index function
	state.SetField(prefix, "__index", state.NewFunction(prefixIndex))
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

func pushPrefix(prefix *msg.Prefix, state *lua.LState) {
	// Create the userdata
	data := state.NewUserData()
	data.Value = prefix

	// Set the metatable
	state.SetMetatable(data, state.GetTypeMetatable("prefix"))

	// Push the prefix
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

func checkPrefix(state *lua.LState) *msg.Prefix {
	data := state.CheckUserData(1)
	if v, ok := data.Value.(*msg.Prefix); ok {
		return v
	}
	state.ArgError(1, "prefix expected")
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
	case "Prefix":
		pushPrefix(msg.Prefix, state)
	case "Tags":
		tags := state.NewTable()
		for key, value := range msg.Tags {
			tags.RawSetString(key, lua.LString(value))
		}
		state.Push(tags)
	default:
		state.Push(lua.LNil)
	}

	return 1
}

func prefixIndex(state *lua.LState) int {
	prefix := checkPrefix(state)
	key := state.CheckString(2)

	switch key {
	case "Name":
		state.Push(lua.LString(prefix.Name))
	case "User":
		state.Push(lua.LString(prefix.User))
	case "Host":
		state.Push(lua.LString(prefix.Host))
	default:
		state.Push(lua.LNil)
	}

	return 1
}
