package geoffrey

import "github.com/yuin/gopher-lua"
import "log"

// Load will load the module to the Lua state
func Load(L *lua.LState) int {
	// Register the module
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push the module
	L.Push(mod)

	return 1
}

// Register will globally register this module
func Register(L *lua.LState) {
	L.RegisterModule("geoffrey", exports)
}

var exports = map[string]lua.LGFunction{
	"add": add,
}

func add(L *lua.LState) int {
	// Get the table
	table := L.ToTable(-1)

	// Get all table rows
	hostname := table.RawGetString("Hostname")
	port := table.RawGetString("Port")
	secure := table.RawGetString("Secure")
	verify := table.RawGetString("InsecureSkipVerify")
	nick := table.RawGetString("Nick")
	user := table.RawGetString("User")
	name := table.RawGetString("Name")

	// Check the types of all parameters
	if hostname.Type() != lua.LTString {
		log.Fatalln("Hostname must be a string")
	}

	if port.Type() != lua.LTNumber {
		log.Fatalln("Port must be a number")
	}

	if secure.Type() != lua.LTBool {
		log.Fatalln("Secure must be a boolean")
	}

	if verify.Type() != lua.LTBool {
		log.Fatalln("InsecureSkipVerify must be a boolean")
	}

	if nick.Type() != lua.LTString {
		log.Fatalln("Nick must be a string")
	}

	if user.Type() != lua.LTString {
		log.Fatalln("User must be a string")
	}

	if name.Type() != lua.LTString {
		log.Fatalln("Name must be a string")
	}

	return 0
}
