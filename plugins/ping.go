package plugins

import (
	"github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
)

func init() {
	bot.RegisterHandler(PingHandler)
}

// PingHandler will respond to ping requests
var PingHandler = bot.Handler{
	Name:        "Ping",
	Description: "Handles ping requests from the server",
	Event:       irc.Ping,
	Run: func(bot *bot.Bot, msg *msg.Message) error {
		bot.Pong(msg.Trailing)
		return nil
	},
}
