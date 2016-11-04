package geoffrey

import (
	log "github.com/Sirupsen/logrus"
	"github.com/yuin/gopher-lua"
)

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
		log.Errorf("hostname must be 'string' but it was '%s'", hostname.Type().String())
	}

	if port.Type() != lua.LTNumber {
		log.Errorf("port must be 'number' but it was '%s'", port.Type().String())
	}

	if secure.Type() != lua.LTBool {
		log.Errorf("secure must be 'boolean' but it was '%s'", secure.Type().String())
	}

	if verify.Type() != lua.LTBool {
		log.Errorf("insecureSkipVerify must be 'boolean' but it was '%s'", verify.Type().String())
	}

	if nick.Type() != lua.LTString {
		log.Errorf("nick must be 'string' but it was '%s'", nick.Type().String())
	}

	if user.Type() != lua.LTString {
		log.Errorf("use must be 'string' but it was '%s'", user.Type().String())
	}

	if name.Type() != lua.LTString {
		log.Errorf("name must be 'string' but it was '%s'", name.Type().String())
	}

	return 0
}
