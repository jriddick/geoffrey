package helper

import (
	log "github.com/Sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

// GetUserData will check for and return LUserData on pos
// n on the Lua stack. If not found or if a different type
// it will return nil.
func GetUserData(n int, state *lua.LState) *lua.LUserData {
	val := state.Get(n)

	if val.Type() != lua.LTUserData {
		log.WithField("file", state.Where(1)).Errorf("Expected userdata but we got '%s'", val.Type().String())
		return nil
	}

	return val.(*lua.LUserData)
}

// GetString will check for and return LString on pos
// n on the Lua stack. It will return nil if not found.
func GetString(n int, state *lua.LState) *string {
	val := state.Get(n)

	if val.Type() != lua.LTString {
		log.WithField("file", state.Where(1)).Errorf("Expected string but we got '%s'", val.Type().String())
		return nil
	}

	res := lua.LVAsString(val)
	return &res
}

// HasArguments will check if there are n amount of arguments
// on the Lua stack. If there is more ore less arguments we return false.
func HasArguments(n int, state *lua.LState) bool {
	if state.GetTop() != n {
		log.WithField("file", state.Where(1)).Errorf("Expected %d arguments but we got '%d'", n, state.GetTop())
		return false
	}

	return true
}
