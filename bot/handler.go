package bot

import (
	"github.com/jriddick/geoffrey/msg"
)

// HandlerFunc is function signature for OnEvent
// handlers.
type HandlerFunc func(*Bot, *msg.Message) error

// Handler is an OnEvent handler
type Handler struct {
	Name        string
	Description string
	Run         HandlerFunc
	Event       string
}

var (
	// Handlers holds list of all registered handlers
	Handlers map[string]map[string]Handler
)

func init() {
	Handlers = make(map[string]map[string]Handler)
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

	return nil
}
