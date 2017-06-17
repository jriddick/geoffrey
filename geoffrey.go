package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/bot"
)

func init() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Set the log level to debug
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infoln("Geoffrey is now running...")

	config := bot.Config{
		Hostname:           "irc.oftc.net",
		Port:               6697,
		Secure:             true,
		InsecureSkipVerify: false,
		Nick:               "geoffrey",
		User:               "geoffrey",
		Name:               "geoffrey-bot",
		Channels:           []string{"#geoffrey-dev"},
		Timeout:            1000,
		TimeoutLimit:       1000,
		ReconnectLimit:     1000,
	}

	// Get the bot manager
	manager := bot.NewManager()

	// Create the first bot
	bot := bot.NewBot(config)

	// Add the bot to the manager
	if err := manager.Add("oftc", bot); err != nil {
		log.Fatalf("Error: %v (%s)", err, "oftc")
	}

	// Listen and run
	if err := manager.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
