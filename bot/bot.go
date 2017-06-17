package bot

import (
	"fmt"
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
	handlers     map[string]Handler
}

// NewBot creates a new bot
func NewBot(config Config) *Bot {
	// Create the bot
	bot := &Bot{
		client: irc.NewIRC(irc.Config{
			Hostname:           config.Hostname,
			Port:               config.Port,
			Secure:             config.Secure,
			InsecureSkipVerify: config.InsecureSkipVerify,
			Timeout:            time.Second * time.Duration(config.Timeout),
			TimeoutLimit:       config.TimeoutLimit,
		}),
		config:       config,
		stop:         make(chan struct{}),
		handlers:     make(map[string]Handler),
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
	for {
		select {
		case <-b.stop:
			// Send quit message
			b.writer <- "QUIT :Closed"

			// Disconnect the client
			b.client.Disconnect()
			break
		case msg := <-b.reader:
			// Log all messages
			log.Debugln(msg.String())

			// Run all handlers
			for _, handler := range b.handlers {
				if msg.Command == handler.Event {
					// Mark start time
					start := time.Now()

					// Execute the handler
					if err := handler.Run(b, msg); err != nil {
						log.Errorf("%s: %v", handler.Name, err)
					}

					// Log the execution time
					log.Infof("Handler '%s' completed in %s", handler.Name, time.Since(start))
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
	b.writer <- "PING :" + message
}

// Pong will send pong to the server
func (b *Bot) Pong(message string) {
	b.writer <- "PONG :" + message
}

// Nick will send the nick command to the server and
// update the stored nick.
func (b *Bot) Nick(nick string) {
	// Set the nick
	b.config.Nick = nick

	// Send the nick
	b.writer <- "NICK " + nick
}

// User will send the user command to the server and
// update the stored name and user
func (b *Bot) User(user, name string) {
	// Set the stored user and name
	b.config.User = user
	b.config.Name = name

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

// AddHandler adds handler to the bot
func (b *Bot) AddHandler(handler Handler) error {
	// Do not add duplicate handlers
	if _, ok := b.handlers[handler.Name]; ok {
		return ErrHandlerExists
	}

	// Add the handler to the bot
	b.handlers[handler.Name] = handler
	return nil
}
