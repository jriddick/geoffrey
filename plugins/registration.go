package plugins

import (
	"github.com/jriddick/geoffrey/bot"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
)

func init() {
	bot.RegisterHandler(RegistrationHandler)
}

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
			bot.Nick(config.Identification.Nick)
			bot.User(config.Identification.User, config.Identification.Name)
		}
		return nil
	},
}
