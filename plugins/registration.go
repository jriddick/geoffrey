package plugins

import (
	"github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
)

// RegistrationHandler will register the bot to the server
var RegistrationHandler = bot.Handler{
	Name:        "Registration",
	Description: "Registers the bot to the IRC server",
	Event:       irc.Notice,
	Run: func(bot *bot.Bot, msg *msg.Message) error {
		if msg.Trailing == "*** Looking up your hostname..." {
			// Get the configuration
			config := bot.Config()

			// Send our registration details
			bot.Nick(config.Nick)
			bot.User(config.User, config.Name)
		}
		return nil
	},
}
