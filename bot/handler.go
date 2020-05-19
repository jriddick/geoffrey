package bot

import (
	"github.com/jriddick/geoffrey/msg"
)

// HandlerFunc is function signature for OnEvent
// handlers.
type HandlerFunc func(*Bot, *msg.Message) (bool, error)

// InitFunc is run once on bot startup to initialize and
// setup the plugin.
type InitFunc func(*Bot) (bool, error)

// Handler is an OnEvent handler
type Handler struct {
	Name        string
	Description string
	Run         HandlerFunc
	Event       string
	Init        InitFunc
}

var (
	// Handlers holds list of all registered handlers
	Handlers map[string]map[string]Handler

	// HandlerList holds a list of all handler
	HandlerList map[string]Handler
)

func init() {
	Handlers = make(map[string]map[string]Handler)
	HandlerList = make(map[string]Handler)
}

// RegisterHandler will register the given handler for use
// with bots
func RegisterHandler(handler Handler) error {
	// Create the event map if it doesn't exist
	if _, ok := Handlers[handler.Event]; !ok {
		Handlers[handler.Event] = make(map[string]Handler)
	}

	// Add the handler if it doesn't exist
	if _, ok := Handlers[handler.Event][handler.Name]; !ok {
		Handlers[handler.Event][handler.Name] = handler
	} else {
		return ErrHandlerExists
	}

	// Append the handler to the list
	if _, ok := HandlerList[handler.Name]; !ok {
		HandlerList[handler.Name] = handler
	}

	return nil
}
