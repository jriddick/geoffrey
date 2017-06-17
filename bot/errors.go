package bot

import "errors"

var (
	// ErrBotExists occurs when the manager already
	// has a bot with the same name.
	ErrBotExists = errors.New("manager: Bot with that name already exists")
	// ErrBotNotFound occurs when no bot is found with
	// the same name.
	ErrBotNotFound = errors.New("manager: Could not find a bot with that name")
	// ErrHandlerExists occurs when you try to add a handler
	// that already exists
	ErrHandlerExists = errors.New("manager: Handler already exists")
)
