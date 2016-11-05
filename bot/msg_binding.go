package bot

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/helper"
	"github.com/jriddick/geoffrey/msg"
	"github.com/yuin/gopher-lua"
)

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
	// Get the logger
	logger := log.WithField("file", state.Where(1))

	// Try to get userdata
	data := helper.GetUserData(1, state)
	if data != nil {
		// Check if userdata is the right type
		if msg, ok := data.Value.(*msg.Message); ok {
			return msg
		}

		logger.Errorf("Expected userdata type msg.Message but we got '%s'", reflect.TypeOf(data.Value).Name())
	}

	return nil
}

func checkPrefix(state *lua.LState) *msg.Prefix {
	// Get the logger
	logger := log.WithField("file", state.Where(1))

	// Try to get userdata
	data := helper.GetUserData(1, state)
	if data != nil {
		// Check if userdata is the right type
		if prefix, ok := data.Value.(*msg.Prefix); ok {
			return prefix
		}

		logger.Errorf("Expected userdata type msg.Prefix but we got '%s'", reflect.TypeOf(data.Value).Name())
	}

	return nil
}

func messageIndex(state *lua.LState) int {
	msg := checkMessage(state)
	key := helper.GetString(2, state)

	if key == nil || msg == nil {
		return 0
	}

	switch *key {
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
	key := helper.GetString(2, state)

	if prefix == nil || key == nil {
		return 0
	}

	switch *key {
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
