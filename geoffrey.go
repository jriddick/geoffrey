package main

import (
	"os"

	"sync"

	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/bot"
)

var (
	signals = make(chan os.Signal, 1)
	wait    sync.WaitGroup
)

func init() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Set the log level to debug
	log.SetLevel(log.DebugLevel)

	// Get signals
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		// Wait until signal
		<-signals

		// Let the program exit
		wait.Done()
	}()
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

	bot := bot.NewBot(config)

	if err := bot.Connect(); err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}

	// Wait until we should exit
	wait.Add(1)
	wait.Wait()

	// Stop the bot
	bot.Close()
}
