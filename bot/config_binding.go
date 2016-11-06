package bot

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/helper"
	"github.com/yuin/gopher-lua"
)

// RegisterConfig will register the Config struct to Lua
func RegisterConfig(state *lua.LState) {
	// Create the Metatable
	meta := state.NewTypeMetatable("config")
	state.SetGlobal("config", meta)

	// Bind our functions
	state.SetField(meta, "__index", state.NewFunction(configIndex))
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

func checkConfig(state *lua.LState) *Config {
	// Get the logger
	logger := log.WithField("file", state.Where(1))

	// Try to get the userdata
	data := helper.GetUserData(1, state)
	if data != nil {
		if v, ok := data.Value.(*Config); ok {
			return v
		}

		logger.Errorf("Expected userdata type Config but we got '%s'", reflect.TypeOf(data.Value).Name())
	}

	return nil
}

func configIndex(state *lua.LState) int {
	config := checkConfig(state)
	key := helper.GetString(2, state)

	if config == nil || key == nil {
		return 0
	}

	switch *key {
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
			channels.RawSetInt(key+1, lua.LString(value))
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
