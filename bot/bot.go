package bot

import (
	"fmt"
	"net"
	"strings"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jriddick/geoffrey/irc"
	"github.com/jriddick/geoffrey/msg"
)

// MessageHandler is the function type for
// message handlers
type MessageHandler func(*Bot, string)

// Bot is the structure for an IRC bot
type Bot struct {
	client       *irc.IRC
	writer       chan<- string
	reader       <-chan *msg.Message
	stop         chan struct{}
	config       Config
	disconnected chan struct{}
}

// NewBot creates a new bot
func NewBot(config Config) *Bot {
	// Create the bot
	bot := &Bot{
		client: irc.NewIRC(irc.Config{
			Hostname:           config.Hostname,
			Port:               config.Port,
			Secure:             config.Secure.Enable,
			InsecureSkipVerify: !config.Secure.Verify,
			Timeout:            time.Millisecond * time.Duration(config.Timings.Timeout),
		}),
		config:       config,
		stop:         make(chan struct{}),
		disconnected: make(chan struct{}),
	}

	return bot
}

// Connect will connect the bot to the server
func (b *Bot) Connect() error {
	// Connect the client
	if err := b.client.Connect(); err != nil {
		return err
	}

	// Get the reader and writer channels
	b.writer = b.client.Writer()
	b.reader = b.client.Reader()

	return nil
}

// Handler will start processing messages
func (b *Bot) Handler() {
	// Handle error messages in another routine
	go func(bot *Bot) {
		for {
			select {
			case <-b.stop:
				break
			case err := <-b.client.Errors():
				// Log the error that we got
				log.Errorf("[geoffrey] %v", err)

				// Check if timeout error occured
				if err, ok := err.(net.Error); ok && err.Timeout() {
					// Reconnect the bot
					if err := b.client.Reconnect(); err != nil {
						log.Fatalf("[geoffrey] %s", err)
					}
				}
			}
		}
	}(b)

	// Ping the server at specific interval
	go func(bot *Bot) {
		ticker := time.NewTicker(120 * time.Second)
		for {
			select {
			case <-b.stop:
				break
			case <-ticker.C:
				b.Ping(fmt.Sprintf("%d", time.Now().UnixNano()))
			}
		}
	}(b)

	// Handle disconnects and regular error messages
	for {
		select {
		case <-b.stop:
			// Disconnect the client
			b.client.Disconnect("Closed")
			break
		case message := <-b.reader:
			// Log all messages
			log.Debugln(message.String())

			// Get all handlers for this event
			if handlers, ok := Handlers[message.Command]; ok {
				// Go through all configured handlers
				for _, name := range b.config.Plugins {
					// Run the handler if we found it
					if handler, ok := handlers[name]; ok {
						go func(bot *Bot, msg *msg.Message, handler Handler) {
							// Mark start time
							start := time.Now()

							// Execute the handler
							if err := handler.Run(b, msg); err != nil {
								log.Errorf("[%s] %v", handler.Name, err)
							}

							// Log the execution time
							log.Infof("Handler '%s' completed in %s", handler.Name, time.Since(start))
						}(b, message, handler)
					}
				}
			}
		}
	}
}

// Send will send the given message to the given receiver
func (b *Bot) Send(recv, msg string) {
	b.writer <- fmt.Sprintf("PRIVMSG %s :%s", recv, msg)
}

// Join will join the given channel
func (b *Bot) Join(channel string) {
	// Make sure we have a hashtag
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}

	// Send the join command
	b.writer <- fmt.Sprintf("JOIN %s", channel)
}

// Ping will send ping to the server
func (b *Bot) Ping(message string) {
	b.writer <- "PING " + message
}

// Pong will send pong to the server
func (b *Bot) Pong(message string) {
	b.writer <- "PONG :" + message
}

// Nick will send the nick command to the server and
// update the stored nick.
func (b *Bot) Nick(nick string) {
	// Set the nick
	b.config.Identification.Nick = nick

	// Send the nick
	b.writer <- "NICK " + nick
}

// User will send the user command to the server and
// update the stored name and user
func (b *Bot) User(user, name string) {
	// Set the stored user and name
	b.config.Identification.User = user
	b.config.Identification.Name = name

	// Send the command
	b.writer <- "USER " + user + " 0 * :" + name
}

// Close will disconnect the bot from the server
func (b *Bot) Close() {
	close(b.stop)
}

// Config returns the configuration
func (b *Bot) Config() Config {
	return b.config
}
