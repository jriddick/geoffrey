package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/bot"
	_ "github.com/jriddick/geoffrey/plugins"
	"github.com/spf13/viper"
)

func init() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Set the log level to debug
	log.SetLevel(log.DebugLevel)

	// Set the configuration information
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Set the default values for configuration
	viper.SetDefault("logs.location", "logs")
	viper.SetDefault("logs.level", "INFO")
}

func main() {
	log.Infoln("[geoffrey] Running")

	// Load the configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Could not read configuration: %s\n", err)
	}

	// Configure the logger level
	switch viper.GetString("logs.level") {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN", "WARNING":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.Fatalf("[geoffrey] Tried to set log level to '%s'", viper.GetString("logs.level"))
	}

	// Create the manager
	manager := bot.NewManager()

	// Get the bot configurations
	var bots []bot.Config

	// Unmarshal the bots
	if err := viper.UnmarshalKey("bots", &bots); err != nil {
		log.Fatalf("Could not read configuration: %s\n", err)
	}

	// Add all bots to the manager
	for _, config := range bots {
		if err := manager.Add(config.BotName, bot.NewBot(config)); err != nil {
			log.Fatalf("[%s] %v", config.BotName, err)
		} else {
			log.Infof("[geoffrey] Added bot '%s'", config.BotName)
		}
	}

	// Make sure we actaully have a bot registered
	if len(bots) < 1 {
		log.Fatalf("[geoffrey] You need a minimum of one configured bot")
	}

	log.Infof("[geoffrey] Started")

	// Listen and run
	if err := manager.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
