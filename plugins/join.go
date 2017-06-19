package plugins

import (
	"github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
)

func init() {
	bot.RegisterHandler(JoinHandler)
}

// JoinHandler will join all configured channels
// after bot has been registered.
var JoinHandler = bot.Handler{
	Name:        "Join",
	Description: "Joins all pre-defined channels after registration",
	Event:       irc.Welcome,
	Run: func(bot *bot.Bot, msg *msg.Message) error {
		// Get the configuration
		config := bot.Config()

		// Join the channels in the configuration
		for _, channel := range config.Channels {
			bot.Join(channel)
		}

		return nil
	},
}
