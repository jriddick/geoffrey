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
